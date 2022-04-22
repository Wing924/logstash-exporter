package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "logstash"
	statsPath = "/_node/stats"
)

var (
	ErrBadStatus = errors.New("bad status code")
)

type Collector struct {
	URI    string
	mutex  sync.RWMutex
	client *http.Client

	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter
	logstashStatus    prometheus.Gauge
	logstashInfo      *prometheus.Desc

	jvm            *jvmCollector
	process        *processCollector
	pipelineConfig *pipelineConfigCollector
	reloadsConfig  *reloadsConfigCollector
	event          *eventCollector
	pipeline       *pipelinesCollector
}

func NewCollector(uri string, timeout time.Duration) (*Collector, error) {
	if strings.HasSuffix(uri, "/") {
		uri = uri[0 : len(uri)-1]
	}
	_, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return &Collector{
		URI:    uri + statsPath,
		client: client,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of logstash successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total logstash scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_json_parse_failures",
			Help:      "Number of errors while parsing JSON.",
		}),
		logstashStatus: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "status",
			Help:      "Was the logstash status: 0 for Green; 1 for Yellow; 2 for Red.",
		}),
		logstashInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "info"),
			"A metric with a constant '1' value labeled by version, http_address, name, id and ephemeral_id from Logstash instance.",
			[]string{"version", "http_address", "name", "id", "ephemeral_id"},
			nil,
		),
		jvm:            newJVMCollector(),
		process:        newProcessCollector(),
		pipelineConfig: newPipelineConfigCollector(),
		reloadsConfig:  newReloadsConfigCollector(),
		event:          newEventCollector(),
		pipeline:       newPipelinesCollector(),
	}, nil
}

// Describe describes all the metrics ever exported by the logstash exporter.
// It implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect fetches the stats from configured logstash and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects.
	defer c.mutex.Unlock()

	c.up.Set(c.scrape(ch))

	ch <- c.up
	ch <- c.totalScrapes
	ch <- c.jsonParseFailures
	ch <- c.logstashStatus
}

func (c *Collector) scrape(ch chan<- prometheus.Metric) (up float64) {
	c.totalScrapes.Inc()

	body, err := c.fetch()
	if err != nil {
		logrus.WithError(err).Warnln("can't scrape logstash", statsPath)
		return 0
	}
	defer body.Close()

	var stats NodeStats
	if err = json.NewDecoder(body).Decode(&stats); err != nil {
		logrus.WithError(err).Warn("can't parse json")
		c.jsonParseFailures.Inc()
		return 0
	}

	c.logstashStatus.Set(c.getStatus(stats))

	ch <- prometheus.MustNewConstMetric(c.logstashInfo, prometheus.GaugeValue, 1.0,
		stats.Version,
		stats.HttpAddress,
		stats.Name,
		stats.ID,
		stats.EphemeralID,
	)

	c.jvm.Collect(stats.JVM, ch)
	c.process.Collect(stats.Process, ch)
	c.pipelineConfig.Collect(stats.Pipeline, ch)
	c.reloadsConfig.Collect(stats.Reloads, ch)
	c.event.Collect(stats.Event, ch)
	c.pipeline.Collect(stats.Pipelines, ch)

	return 1
}

func (c *Collector) getStatus(stats NodeStats) float64 {
	switch stats.Status {
	case "green":
		return 0
	case "yellow":
		return 1
	default:
		return 2
	}
}

func (c *Collector) fetch() (io.ReadCloser, error) {
	resp, err := c.client.Get(c.URI)
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP status %d: %w", resp.StatusCode, ErrBadStatus)
	}
	return resp.Body, nil
}
