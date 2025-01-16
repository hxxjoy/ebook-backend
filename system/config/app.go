package config

import "time"

type AppConfig struct {
	Name        string        `yaml:"name"`
	Mode        string        `yaml:"mode"` // dev, test, prod
	Version     string        `yaml:"version"`
	JWTSecret   string        `yaml:"jwt_secret"`
	JWTTimeout  time.Duration `yaml:"jwt_timeout"`
	PageSize    int           `yaml:"page_size"`
	MaxFileSize int           `yaml:"max_file_size"` // 单位 MB
}
