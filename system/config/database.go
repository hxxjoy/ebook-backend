package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Driver          string        `yaml:"driver"` // mysql, postgres
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.Username, c.Password, c.Database)
	default:
		return ""
	}
}
