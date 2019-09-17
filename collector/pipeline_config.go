package collector

import "github.com/prometheus/client_golang/prometheus"

type pipelineConfigCollector struct {
	Workers    *prometheus.Desc
	BatchSize  *prometheus.Desc
	BatchDelay *prometheus.Desc
}

func newPipelineConfigCollector() *pipelineConfigCollector {
	desc := newDescFunc(namespace, "pipeline_config")
	return &pipelineConfigCollector{
		Workers:    desc("workers", "The number of workers that will, in parallel, execute the filter and output stages of the pipeline."),
		BatchSize:  desc("batch_size", "The maximum number of events an individual worker thread will collect from inputs before attempting to execute its filters and outputs."),
		BatchDelay: desc("batch_delay_seconds", "How long to wait before dispatching an undersized batch to workers."),
	}
}

func (c *pipelineConfigCollector) Collect(p PipelineConfig, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.Workers, prometheus.GaugeValue, float64(p.Workers))
	ch <- prometheus.MustNewConstMetric(c.BatchSize, prometheus.GaugeValue, float64(p.BatchSize))
	ch <- prometheus.MustNewConstMetric(c.BatchDelay, prometheus.GaugeValue, float64(p.BatchDelay)/1000.0)
}
