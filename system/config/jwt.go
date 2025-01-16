package config

import "time"

type JWTConfig struct {
    AccessTokenSecret  string        `yaml:"access_token_secret"`
    RefreshTokenSecret string        `yaml:"refresh_token_secret"`
    AccessTokenExpire  time.Duration  `yaml:"access_token_expire"` // 短期，如 15分钟
    RefreshTokenExpire time.Duration  `yaml:"refresh_token_expire"`// 长期，如 30天
}