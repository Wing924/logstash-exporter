package collector

import "github.com/prometheus/client_golang/prometheus"

type reloadsConfigCollector struct {
	Failures  *prometheus.Desc
	Successes *prometheus.Desc
}

func newReloadsConfigCollector() *reloadsConfigCollector {
	desc := newDescFunc(namespace, "reloads_config")
	return &reloadsConfigCollector{
		Failures:  desc("failures_total", "Number of failures during config reload."),
		Successes: desc("successes_total", "Number of successful config reloads."),
	}
}

func (c *reloadsConfigCollector) Collect(p ReloadsConfig, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.Failures, prometheus.CounterValue, float64(p.Failures))
	ch <- prometheus.MustNewConstMetric(c.Successes, prometheus.CounterValue, float64(p.Successes))
}
