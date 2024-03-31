package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Environment string   `envconfig:"ENV" default:"local"`
	AppName     string   `envconfig:"APP_NAME" default:"payment"`
	LogLevel    string   `envconfig:"LOG_LEVEL" default:"info"`
	Port        string   `envconfig:"PORT" default:"8081"`
	SecretKey   string   `envconfig:"SECRET_KEY" default:"someUltraSecretKey"`
	Database    DBConfig `envconfig:"DATABASE"`
}

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	Username string `envconfig:"DB_USERNAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

var c Configuration

func Config() Configuration {
	return c
}

func Environment() {
	if err := envconfig.Process("", &c); err != nil {
		if err != nil {
			fmt.Printf("Error loading configuration: %s\n", err)
		}
		return
	}
}
