package config

import (
	"time"
)

type Config struct {
	HTTP     HTTP     `envPrefix:"HTTP_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
	Logger   Logger   `envPrefix:"LOGGER_"`
	JWT      JWT      `envPrefix:"JWT_"`
	SMTP     SMTP     `envPrefix:"SMTP_"`
}

type HTTP struct {
	Url             string        `env:"URL"`
	Port            int           `env:"PORT"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

type Postgres struct {
	Host            string        `env:"HOST"`
	Port            int           `env:"PORT"`
	Username        string        `env:"USERNAME"`
	Password        string        `env:"PASSWORD"`
	Dbname          string        `env:"DBNAME"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
}

type Logger struct {
	LogLevel    string `env:"LOG_LEVEL"`
	FileLogName string `env:"FILE_LOG_NAME"`
	MaxBackups  int    `env:"MAX_BACKUPS"`
	MaxAge      int    `env:"MAX_AGE"`
	MaxSize     int    `env:"MAX_SIZE"`
	Compress    bool   `env:"COMPRESS"`
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
