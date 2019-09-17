package collector

import "github.com/prometheus/client_golang/prometheus"

type eventCollector struct {
	In                *prometheus.Desc
	Filtered          *prometheus.Desc
	Out               *prometheus.Desc
	Duration          *prometheus.Desc
	QueuePushDuration *prometheus.Desc
}

func newEventCollector() *eventCollector {
	desc := newDescFunc(namespace, "event")
	return &eventCollector{
		In:                desc("in_total", "The total number of events in."),
		Filtered:          desc("filtered_total", "The total numbers of filtered."),
		Out:               desc("out_total", "The total number of events out."),
		Duration:          desc("duration_seconds_total", "The total process duration time in seconds."),
		QueuePushDuration: desc("queue_push_duration_seconds_total", "The total in queue duration time in seconds."),
	}
}

func (c *eventCollector) Collect(e Event, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.In, prometheus.CounterValue, float64(e.In))
	ch <- prometheus.MustNewConstMetric(c.Filtered, prometheus.CounterValue, float64(e.Filtered))
	ch <- prometheus.MustNewConstMetric(c.Out, prometheus.CounterValue, float64(e.Out))

	ch <- prometheus.MustNewConstMetric(
		c.Duration, prometheus.CounterValue, float64(e.DurationInMillis)/1000.0)
	ch <- prometheus.MustNewConstMetric(
		c.QueuePushDuration, prometheus.CounterValue, float64(e.QueuePushDurationInMillis)/1000.0)
}
