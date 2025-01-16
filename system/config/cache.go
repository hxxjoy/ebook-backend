package config

type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
	TTL     int  `yaml:"ttl"`
	Clearup int  `yaml:"clearup"`
	MaxSize int64 `yaml:"max_size"`
}
