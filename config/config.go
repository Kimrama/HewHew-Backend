package config

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *Server   `mapstructure:"server" validate:"required"`
		Auth     *Auth     `mapstructure:"auth" validate:"required"`
		Database *Database `mapstructure:"database" validate:"required"`
	}

	Server struct {
		Port           int      `mapstructure:"port" validate:"required"`
		AllowedOrigins []string `mapstructure:"allowed_origins" validate:"required"`
	}
	Auth struct {
		Secret string `mapstructure:"secret" validate:"required"`
	}
	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		DBName   string `mapstructure:"dbname" validate:"required"`
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}

		validating := validator.New()
		if err := validating.Struct(configInstance); err != nil {
			panic(err)
		}

	})
	return configInstance
}
