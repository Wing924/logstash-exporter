# logstash-exporter
Prometheus exporter for the metrics available in Logstash since version 7.0.

## Usage

### Install

Using docker
```bash
docker pull wing924/logstash-exporter
docker run wing924/logstash-exporter
```

Build from source
```bash
git clone https://github.com/Wing924/logstash-exporter
cd logstash-exporter
make
```

### Flags

```bash
logstash-exporter --help
usage: logstash-exporter [<flags>]

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9649"
                             Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                             Path under which to expose metrics.
      --logstash.scrape-uri="http://localhost:9600"
                             URI on which to scrape logstash.
      --logstash.timeout=5s  Timeout for trying to get stats from logstash.
      --version              Show application version.
```

## Implemented Metrics

* metadata/config metrics
  * `logstash_exporter_build_info` A metric with a constant '1' value labeled by version, revision, branch, and goversion from which logstash_exporter was built.
  * `logstash_exporter_json_parse_failures` Number of errors while parsing JSON.
  * `logstash_exporter_total_scrapes` Current total logstash scrapes.
  * `logstash_info` A metric with a constant '1' value labeled by version, http_address, name, id and ephemeral_id from Logstash instance.
  * `logstash_pipeline_config_batch_delay_seconds` How long to wait before dispatching an undersized batch to workers.
  * `logstash_pipeline_config_batch_size` The maximum number of events an individual worker thread will collect from inputs before attempting to execute its filters and outputs.
  * `logstash_pipeline_config_workers` The number of workers that will, in parallel, execute the filter and output stages of the pipeline.
  * `logstash_up` Was the last scrape of logstash successful.
* event metrics
  * `logstash_event_duration_seconds_total` The total process duration time in seconds.
  * `logstash_event_filtered_total` The total numbers of filtered.
  * `logstash_event_in_total` The total number of events in.
  * `logstash_event_out_total` The total number of events out.
  * `logstash_event_queue_push_duration_seconds_total` The total in queue duration time in seconds.
* JVM metrics
  * `logstash_jvm_gc_collection_duration_seconds` GC collection duration.
  * `logstash_jvm_heap_committed_bytes` Current JVM heap committed size
  * `logstash_jvm_heap_used_bytes` Current JVM heap used size
  * `logstash_jvm_heap_used_ratio` Current JVM heap usage ratio.
  * `logstash_jvm_memory_pool_committed_bytes` Current JVM heap pool committed size
  * `logstash_jvm_memory_pool_max_bytes` Current JVM heap pool max size
  * `logstash_jvm_memory_pool_used_bytes` Current JVM heap pool used size
  * `logstash_jvm_threads_count` Current JVM thread count.
* pipeline metrics
  * `logstash_pipeline_event_duration_seconds_total` The total process duration time in seconds.
  * `logstash_pipeline_event_filtered_total` The total numbers of filtered.
  * `logstash_pipeline_event_in_total` The total number of events in.
  * `logstash_pipeline_event_out_total` The total number of events out.
  * `logstash_pipeline_event_queue_push_duration_seconds_total` The total in queue duration time in seconds.
  * `logstash_pipeline_filter_duration_seconds_total` The total process duration time in seconds
  * `logstash_pipeline_filter_in_total` The total number of events in.
  * `logstash_pipeline_filter_out_total` The total number of events out.
  * `logstash_pipeline_input_connections` The current number of connections.
  * `logstash_pipeline_input_out_total` The total number of events out.
  * `logstash_pipeline_input_queue_push_seconds_total` The total in queue duration time in seconds
  * `logstash_pipeline_output_duration_seconds_total` The total process duration time in seconds
  * `logstash_pipeline_output_in_total` The total number of events in.
  * `logstash_pipeline_output_out_total` The total number of events out.
  * `logstash_pipeline_queue_event_count` The current events in queue.
  * `logstash_pipeline_queue_max_size_bytes` The max queue size in bytes.
  * `logstash_pipeline_queue_size_bytes` The current queue size in bytes.
* process metrics
  * `logstash_process_cpu_usage_ratio` Was the CPU usage
  * `logstash_process_load_average` Was the system load average
  * `logstash_process_max_file_descriptors` Max file descriptors
  * `logstash_process_open_file_descriptors` Current open file descriptors
  * `logstash_process_process_time_seconds` Was the total process time.
  * `logstash_process_total_virtual_memory_bytes` Was the used virtual memory.
  * `logstash_status` Was the logstash status: 0 for Green; 1 for Yellow; 2 for Red.
* reloads metrics
  * `logstash_reloads_config_failures_total` Number of failures during config reload
  * `logstash_reloads_config_successes_total` Number of successful config reloads
