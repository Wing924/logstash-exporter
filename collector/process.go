package collector

import "github.com/prometheus/client_golang/prometheus"

type processCollector struct {
	openFileDescriptors *prometheus.Desc
	maxFileDescriptors  *prometheus.Desc
	totalVirtualMemory  *prometheus.Desc
	processTime         *prometheus.Desc
	cpuUsage            *prometheus.Desc
	loadAverage         *prometheus.Desc
}

func newProcessCollector() *processCollector {
	desc := newDescFunc(namespace, "process")
	return &processCollector{
		openFileDescriptors: desc("open_file_descriptors", "Current open file descriptors"),
		maxFileDescriptors:  desc("max_file_descriptors", "Max file descriptors"),
		totalVirtualMemory:  desc("total_virtual_memory_bytes", "Was the used virtual memory."),
		processTime:         desc("process_time_seconds", "Was the total process time."),
		cpuUsage:            desc("cpu_usage_ratio", "Was the CPU usage"),
		loadAverage:         desc("load_average", "Was the system load average", "load"),
	}
}

func (c *processCollector) Collect(p Process, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.openFileDescriptors, prometheus.GaugeValue, float64(p.OpenFileDescriptors))
	ch <- prometheus.MustNewConstMetric(c.maxFileDescriptors, prometheus.GaugeValue, float64(p.MaxFileDescriptors))
	ch <- prometheus.MustNewConstMetric(c.totalVirtualMemory, prometheus.GaugeValue, float64(p.Mem.TotalVirtualInBytes))
	ch <- prometheus.MustNewConstMetric(c.processTime, prometheus.CounterValue, float64(p.CPU.TotalInMillis)/1000.0)
	ch <- prometheus.MustNewConstMetric(c.cpuUsage, prometheus.GaugeValue, float64(p.CPU.Percent)/100.0)
	ch <- prometheus.MustNewConstMetric(c.loadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load1, "1")
	ch <- prometheus.MustNewConstMetric(c.loadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load5, "5")
	ch <- prometheus.MustNewConstMetric(c.loadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load15, "15")
}
