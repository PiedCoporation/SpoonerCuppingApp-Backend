package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTP HTTP `envPrefix:"HTTP_"`
	JWT  JWT  `envPrefix:"JWT_"`
	SMTP SMTP `envPrefix:"SMTP_"`
}

type HTTP struct {
	Url  string `env:"URL"`
	Port int    `env:"PORT"`
}

type JWT struct {
	AccessTokenKey       string        `env:"ACCESS_TOKEN_KEY"`
	AccessTokenExpiresIn time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN"`

	RefreshTokenKey       string        `env:"REFRESH_TOKEN_KEY"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN"`

	RegisterTokenKey       string        `env:"REGISTER_TOKEN_KEY"`
	RegisterTokenExpiresIn time.Duration `env:"REGISTER_TOKEN_EXPIRES_IN"`

	LoginTokenKey       string        `env:"LOGIN_TOKEN_KEY"`
	LoginTokenExpiresIn time.Duration `env:"LOGIN_TOKEN_EXPIRES_IN"`
}

type SMTP struct {
	Host        string `env:"HOST"`
	Port        int    `env:"PORT"`
	Username    string `env:"USERNAME"`
	AppPassword string `env:"APP_PASSWORD"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
