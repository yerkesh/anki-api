package config

import (
	"fmt"
	"net/url"
)

type Config struct {
	DatabaseConfig   DatabaseConfig
	StorageAccessKey string `env:"STORAGE_ACCESS_KEY"`
	StorageSecretKey string `env:"STORAGE_SECRET_KEY"`
	ChatGPTKey       string `env:"CHAT_GPT_API_KEY" env-required:"true"`
	HTTP             HTTP
}

type HTTP struct {
	Port string `env:"PORT" env-default:"8080"`
}

type DatabaseConfig struct {
	Host         string `env:"DATABASE_HOST" env-required:"true"`
	Port         string `env:"DATABASE_PORT" env-required:"true"`
	User         string `env:"DATABASE_USER" env-required:"true"`
	Password     string `env:"DATABASE_PASSWORD" env-required:"true"`
	Name         string `env:"DATABASE_NAME" env-required:"true"`
	PoolMaxConns int    `env:"DATABASE_MAX_CONNECTIONS" env-required:"true"`
}

func (d DatabaseConfig) DSN() string {
	// Use url.URL so that weird chars are escaped automatically.
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%s", d.Host, d.Port),
		Path:   d.Name,
	}

	q := url.Values{}

	//if d.PoolMaxConns > 0 {
	//	q.Set("pool_max_conns", fmt.Sprintf("%d", d.PoolMaxConns))
	//}
	u.RawQuery = q.Encode()

	return u.String()
}

func (d DatabaseConfig) GetHost() string {
	return d.Host
}

func (d DatabaseConfig) GetPort() string {
	return d.Port
}
