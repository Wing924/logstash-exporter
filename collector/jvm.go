package collector

import "github.com/prometheus/client_golang/prometheus"

type jvmCollector struct {
	threadsCount         *prometheus.Desc
	heapUsedRatio        *prometheus.Desc
	heapCommittedInBytes *prometheus.Desc
	heapUsedInBytes      *prometheus.Desc
	poolUsedBytes        *prometheus.Desc
	poolCommittedBytes   *prometheus.Desc
	poolMaxBytes         *prometheus.Desc
	gc                   *prometheus.Desc
}

func newJVMCollector() *jvmCollector {
	desc := newDescFunc(namespace, "jvm")
	return &jvmCollector{
		threadsCount:         desc("threads_count", "Current JVM thread count."),
		heapUsedRatio:        desc("heap_used_ratio", "Current JVM heap usage ratio."),
		heapCommittedInBytes: desc("heap_committed_bytes", "Current JVM heap committed size"),
		heapUsedInBytes:      desc("heap_used_bytes", "Current JVM heap used size"),
		poolUsedBytes:        desc("memory_pool_used_bytes", "Current JVM heap pool used size", "pool"),
		poolCommittedBytes:   desc("memory_pool_committed_bytes", "Current JVM heap pool committed size", "pool"),
		poolMaxBytes:         desc("memory_pool_max_bytes", "Current JVM heap pool max size", "pool"),
		gc:                   desc("gc_collection_duration_seconds", "GC collection duration.", "collector"),
	}
}

func (c *jvmCollector) Collect(jvm JVM, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.threadsCount, prometheus.GaugeValue, float64(jvm.Threads.Count))

	ch <- prometheus.MustNewConstMetric(c.heapUsedRatio, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedPercent)/100.0)
	ch <- prometheus.MustNewConstMetric(c.heapCommittedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapCommittedInBytes))
	ch <- prometheus.MustNewConstMetric(c.heapUsedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedInBytes))

	ch <- prometheus.MustNewConstMetric(c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.UsedInBytes), "young")
	ch <- prometheus.MustNewConstMetric(c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.UsedInBytes), "survivor")
	ch <- prometheus.MustNewConstMetric(c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.UsedInBytes), "old")

	ch <- prometheus.MustNewConstMetric(c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.CommittedInBytes), "young")
	ch <- prometheus.MustNewConstMetric(c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.CommittedInBytes), "survivor")
	ch <- prometheus.MustNewConstMetric(c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.CommittedInBytes), "old")

	ch <- prometheus.MustNewConstMetric(c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.MaxInBytes), "young")
	ch <- prometheus.MustNewConstMetric(c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.MaxInBytes), "survivor")
	ch <- prometheus.MustNewConstMetric(c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.MaxInBytes), "old")

	ch <- prometheus.MustNewConstSummary(c.gc, jvm.GC.Collectors.Young.CollectionCount, float64(jvm.GC.Collectors.Young.CollectionTimeInMillis)/1000.0, nil, "young")
	ch <- prometheus.MustNewConstSummary(c.gc, jvm.GC.Collectors.Old.CollectionCount, float64(jvm.GC.Collectors.Old.CollectionTimeInMillis)/1000.0, nil, "old")
}
