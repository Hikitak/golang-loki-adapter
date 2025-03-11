package models

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	QueueTable   string     `yaml:"queue_table"`
	PollInterval int        `yaml:"poll_interval"`
	Loki         LokiConfig `yaml:"loki"`
}

type LokiConfig struct {
	URL       string            `yaml:"url"`
	Labels    map[string]string `yaml:"labels"`
	Timeout   int               `yaml:"timeout"`
	Retries   int               `yaml:"retries"`
	BatchSize int               `yaml:"batch_size"`
}

type QueueRecord struct {
	ID        int       `json:"id"`
	Data      string    `json:"data"`
}
