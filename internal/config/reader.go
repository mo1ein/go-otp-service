package config

import (
	"errors"

	"fmt"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	httpCfg := HTTP{
		APIHost: loadString("API_HTTP_HOST"),
		APIPort: loadInt("API_HTTP_PORT"),
	}

	postgresCfg := Postgres{
		Host:         loadString("DATABASE_POSTGRES_HOST"),
		Port:         loadInt("DATABASE_POSTGRES_PORT"),
		User:         loadString("DATABASE_POSTGRES_USER"),
		Password:     loadString("DATABASE_POSTGRES_PASSWORD"),
		DatabaseName: loadString("DATABASE_POSTGRES_NAME"),
	}

	redisCfg := Redis{
		Host:     loadString("DATABASE_REDIS_HOST"),
		Port:     loadString("DATABASE_REDIS_PORT"),
		Password: loadString("DATABASE_REDIS_PASSWORD"),
		Database: loadInt("DATABASE_REDIS_DATABASE"),
	}

	return &Config{
		HTTP: httpCfg,
		Database: Database{
			Postgres: postgresCfg,
			Redis:    redisCfg,
		},
	}, nil
}
