package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type pipelinesCollector struct {
	// Event
	In                *prometheus.Desc
	Filtered          *prometheus.Desc
	Out               *prometheus.Desc
	Duration          *prometheus.Desc
	QueuePushDuration *prometheus.Desc

	// Input Plugins
	InputConnections       *prometheus.Desc
	InputQueuePushDuration *prometheus.Desc
	InputOut               *prometheus.Desc

	// Filter Plugins
	FilterDuration *prometheus.Desc
	FilterIn       *prometheus.Desc
	FilterOut      *prometheus.Desc

	// Output Plugins
	OutputDuration *prometheus.Desc
	OutputIn       *prometheus.Desc
	OutputOut      *prometheus.Desc

	// Queue
	EventsCount  *prometheus.Desc
	QueueSize    *prometheus.Desc
	MaxQueueSize *prometheus.Desc
}

func newPipelinesCollector() *pipelinesCollector {
	desc := newDescFunc(namespace, "pipeline")
	return &pipelinesCollector{
		In:                desc("event_in_total", "The total number of events in.", "pipeline"),
		Filtered:          desc("event_filtered_total", "The total numbers of filtered.", "pipeline"),
		Out:               desc("event_out_total", "The total number of events out.", "pipeline"),
		Duration:          desc("event_duration_seconds_total", "The total process duration time in seconds.", "pipeline"),
		QueuePushDuration: desc("event_queue_push_duration_seconds_total", "The total in queue duration time in seconds.", "pipeline"),

		InputConnections:       desc("input_connections", "The current number of connections.", "pipeline", "id", "name"),
		InputQueuePushDuration: desc("input_queue_push_seconds_total", "The total in queue duration time in seconds", "pipeline", "id", "name"),
		InputOut:               desc("input_out_total", "The total number of events out.", "pipeline", "id", "name"),

		FilterDuration: desc("filter_duration_seconds_total", "The total process duration time in seconds", "pipeline", "id", "name", "index"),
		FilterIn:       desc("filter_in_total", "The total number of events in.", "pipeline", "id", "name", "index"),
		FilterOut:      desc("filter_out_total", "The total number of events out.", "pipeline", "id", "name", "index"),

		OutputDuration: desc("output_duration_seconds_total", "The total process duration time in seconds", "pipeline", "id", "name"),
		OutputIn:       desc("output_in_total", "The total number of events in.", "pipeline", "id", "name"),
		OutputOut:      desc("output_out_total", "The total number of events out.", "pipeline", "id", "name"),

		EventsCount:  desc("queue_event_count", "The current events in queue.", "pipeline", "queue_type"),
		QueueSize:    desc("queue_size_bytes", "The current queue size in bytes.", "pipeline", "queue_type"),
		MaxQueueSize: desc("queue_max_size_bytes", "The max queue size in bytes.", "pipeline", "queue_type"),
	}
}

func (c *pipelinesCollector) Collect(p map[string]Pipeline, ch chan<- prometheus.Metric) {
	for pipelineName, pipeline := range p {
		c.collectEvent(pipelineName, pipeline, ch)
		c.collectQueue(pipelineName, pipeline, ch)
		for _, plugin := range pipeline.Plugins.Inputs {
			c.collectInput(pipelineName, plugin, ch)
		}
		for idx, plugin := range pipeline.Plugins.Filters {
			c.collectFilter(pipelineName, idx, plugin, ch)
		}
		for _, plugin := range pipeline.Plugins.Outputs {
			c.collectOutput(pipelineName, plugin, ch)
		}
	}
}

func (c *pipelinesCollector) collectEvent(pipelineName string, p Pipeline, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.In, prometheus.CounterValue, float64(p.Event.In), pipelineName)
	ch <- prometheus.MustNewConstMetric(c.Filtered, prometheus.CounterValue, float64(p.Event.Filtered), pipelineName)
	ch <- prometheus.MustNewConstMetric(c.Out, prometheus.CounterValue, float64(p.Event.Out), pipelineName)

	ch <- prometheus.MustNewConstMetric(
		c.Duration, prometheus.CounterValue, float64(p.Event.DurationInMillis)/1000.0, pipelineName)
	ch <- prometheus.MustNewConstMetric(
		c.QueuePushDuration, prometheus.CounterValue, float64(p.Event.QueuePushDurationInMillis)/1000.0, pipelineName)
}

func (c *pipelinesCollector) collectQueue(pipelineName string, p Pipeline, ch chan<- prometheus.Metric) {
	queueType := p.Queue.Type
	if queueType == "" {
		return
	}
	ch <- prometheus.MustNewConstMetric(c.EventsCount, prometheus.GaugeValue, float64(p.Queue.EventsCount), pipelineName, queueType)
	ch <- prometheus.MustNewConstMetric(c.QueueSize, prometheus.CounterValue, float64(p.Queue.QueueSizeInBytes), pipelineName, queueType)
	ch <- prometheus.MustNewConstMetric(c.MaxQueueSize, prometheus.CounterValue, float64(p.Queue.MaxQueueSizeInBytes), pipelineName, queueType)
}

func (c *pipelinesCollector) collectInput(pipelineName string, p InputPlugin, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.InputConnections, prometheus.GaugeValue, float64(p.CurrentConnections), pipelineName, p.ID, p.Name)
	ch <- prometheus.MustNewConstMetric(
		c.InputQueuePushDuration, prometheus.CounterValue, float64(p.Events.QueuePushDurationInMillis)/1000.0, pipelineName, p.ID, p.Name)
	ch <- prometheus.MustNewConstMetric(c.InputOut, prometheus.CounterValue, float64(p.Events.Out), pipelineName, p.ID, p.Name)
}

func (c *pipelinesCollector) collectFilter(pipelineName string, index int, p FilterPlugin, ch chan<- prometheus.Metric) {
	idx := strconv.Itoa(index)
	ch <- prometheus.MustNewConstMetric(
		c.FilterDuration, prometheus.CounterValue, float64(p.Events.DurationInMillis)/1000.0, pipelineName, p.ID, p.Name, idx)
	ch <- prometheus.MustNewConstMetric(c.FilterIn, prometheus.CounterValue, float64(p.Events.In), pipelineName, p.ID, p.Name, idx)
	ch <- prometheus.MustNewConstMetric(c.FilterOut, prometheus.CounterValue, float64(p.Events.Out), pipelineName, p.ID, p.Name, idx)
}

func (c *pipelinesCollector) collectOutput(pipelineName string, p OutputPlugin, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.OutputDuration, prometheus.CounterValue, float64(p.Events.DurationInMillis)/1000.0, pipelineName, p.ID, p.Name)
	ch <- prometheus.MustNewConstMetric(c.OutputIn, prometheus.CounterValue, float64(p.Events.In), pipelineName, p.ID, p.Name)
	ch <- prometheus.MustNewConstMetric(c.OutputOut, prometheus.CounterValue, float64(p.Events.Out), pipelineName, p.ID, p.Name)
}
