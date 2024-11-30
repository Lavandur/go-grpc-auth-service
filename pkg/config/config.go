package config

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	App `mapstructure:"app"`
}

type (
	App struct {
		Auth            Auth   `mapstructure:"auth"`
		DB              DB     `mapstructure:"db"`
		Level           string `mapstructure:"log-level"`
		DefaultRoleName string `mapstructure:"default-role"`
	}
	Auth struct {
		AccessTokenLifeTime  int `mapstructure:"access-life-time"`
		RefreshTokenLifeTime int `mapstructure:"refresh-life-time"`
	}
	DB struct {
		PostgresQL `mapstructure:"postgresql"`
	}
	PostgresQL struct {
		PostgresqlHost     string `mapstructure:"host"`
		PostgresqlPort     string `mapstructure:"port"`
		PostgresqlUser     string `mapstructure:"user"`
		PostgresqlPassword string `mapstructure:"password"`
		PostgresqlDatabase string `mapstructure:"database"`
		PostgresqlSSLMode  string `mapstructure:"sslmode"`
	}
)

func configureConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logrus.WithField(
				"config file",
				viper.ConfigFileUsed()).Fatal("Config file not found")
		}
		return err
	}
	return nil
}

func SetupConfiguration() (*Config, error) {
	once.Do(func() {
		err := configureConfig()
		if err != nil {
			logrus.Error(err)
		}

		err = viper.Unmarshal(&instance)
		if err != nil {
			logrus.Error(err)
		}
	})
	return instance, nil
}
