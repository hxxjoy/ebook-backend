package config

type LogConfig struct {
	Level      string `yaml:"level"` // debug, info, warn, error
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`    // 单位 MB
	MaxBackups int    `yaml:"max_backups"` // 最大备份数
	MaxAge     int    `yaml:"max_age"`     // 最大保留天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}
