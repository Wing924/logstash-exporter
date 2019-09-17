package main

import (
	"net/http"

	"github.com/Wing924/logstash-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		listenAddress     = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9649").String()
		metricsPath       = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		logstashScrapeURI = kingpin.Flag("logstash.scrape-uri", "URI on which to scrape logstash.").Default("http://localhost:9600").String()
		logstashTimeout   = kingpin.Flag("logstash.timeout", "Timeout for trying to get stats from logstash.").Default("5s").Duration()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("logstash-exporter"))
	kingpin.Parse()

	logrus.WithFields(logrus.Fields{
		"version": version.Info(),
		"build":   version.BuildContext(),
	}).Info("Starting logstash-exporter")

	exporter, err := collector.NewCollector(*logstashScrapeURI, *logstashTimeout)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create exporter")
	}
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("logstash_exporter"))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>Logstash Collector</title></head>
             <body>
             <h1>Logstash Collector</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	logrus.WithField("address", *listenAddress).Info("listening...")
	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}
