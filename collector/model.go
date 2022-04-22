package collector

type NodeStats struct {
	// Top level
	Host        string `json:"host"`
	Version     string `json:"version"`
	HttpAddress string `json:"http_address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	EphemeralID string `json:"ephemeral_id"`
	Status      string `json:"status"`

	Pipeline  PipelineConfig      `json:"pipeline"`
	Reloads   ReloadsConfig       `json:"reloads"`
	JVM       JVM                 `json:"jvm"`
	Process   Process             `json:"process"`
	Event     Event               `json:"events"`
	Pipelines map[string]Pipeline `json:"pipelines"`
}

type PipelineConfig struct {
	Workers    int `json:"workers"`
	BatchSize  int `json:"batch_size"`
	BatchDelay int `json:"batch_delay"`
}

type ReloadsConfig struct {
	Failures  int `json:"failures"`
	Successes int `json:"successes"`
}

type Process struct {
	OpenFileDescriptors     int `json:"open_file_descriptors"`
	PeakOpenFileDescriptors int `json:"peak_open_file_descriptors"`
	MaxFileDescriptors      int `json:"max_file_descriptors"`
	Mem                     struct {
		TotalVirtualInBytes int `json:"total_virtual_in_bytes"`
	} `json:"mem"`
	CPU struct {
		TotalInMillis int `json:"total_in_millis"`
		Percent       int `json:"percent"`
		LoadAverage   struct {
			Load1  float64 `json:"1m"`
			Load5  float64 `json:"5m"`
			Load15 float64 `json:"15m"`
		} `json:"load_average"`
	} `json:"cpu"`
}

type JVM struct {
	Threads struct {
		Count int `json:"count"`
	} `json:"threads"`
	Mem struct {
		HeapUsedPercent      int `json:"heap_used_percent"`
		HeapCommittedInBytes int `json:"heap_committed_in_bytes"`
		HeapUsedInBytes      int `json:"heap_used_in_bytes"`
		Pools                struct {
			Survivor JvmPool `json:"survivor"`
			Old      JvmPool `json:"old"`
			Young    JvmPool `json:"young"`
		} `json:"pools"`
	} `json:"mem"`
	GC struct {
		Collectors struct {
			Old   GCCollector `json:"old"`
			Young GCCollector `json:"young"`
		} `json:"collectors"`
	} `json:"gc"`
	UptimeInMillis int `json:"uptime_in_millis"`
}

type Event struct {
	In                        int `json:"in"`
	Filtered                  int `json:"filtered"`
	Out                       int `json:"out"`
	DurationInMillis          int `json:"duration_in_millis"`
	QueuePushDurationInMillis int `json:"queue_push_duration_in_millis"`
}

type JvmPool struct {
	PeakUsedInBytes  int `json:"peak_used_in_bytes"`
	UsedInBytes      int `json:"used_in_bytes"`
	CommittedInBytes int `json:"committed_in_bytes"`
	PeakMaxInBytes   int `json:"peak_max_in_bytes"`
	MaxInBytes       int `json:"max_in_bytes"`
}

type GCCollector struct {
	CollectionTimeInMillis int    `json:"collection_time_in_millis"`
	CollectionCount        uint64 `json:"collection_count"`
}

type Pipeline struct {
	Event   Event `json:"events"`
	Plugins struct {
		Inputs  []InputPlugin  `json:"inputs"`
		Filters []FilterPlugin `json:"filters"`
		Outputs []OutputPlugin `json:"outputs"`
	} `json:"plugins"`
	Queue struct {
		Type                string `json:"type"`
		EventsCount         int    `json:"events_count"`
		QueueSizeInBytes    int    `json:"queue_size_in_bytes"`
		MaxQueueSizeInBytes int    `json:"max_queue_size_in_bytes"`
	} `json:"queue"`
}

type InputPlugin struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	CurrentConnections int    `json:"current_connections"`
	Events             struct {
		QueuePushDurationInMillis int `json:"queue_push_duration_in_millis"`
		Out                       int `json:"out"`
	} `json:"events"`
}

type FilterPlugin struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Events struct {
		In               int `json:"in"`
		DurationInMillis int `json:"duration_in_millis"`
		Out              int `json:"out"`
	} `json:"events"`
}

type OutputPlugin struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Events struct {
		In               int `json:"in"`
		DurationInMillis int `json:"duration_in_millis"`
		Out              int `json:"out"`
	} `json:"events"`
}
