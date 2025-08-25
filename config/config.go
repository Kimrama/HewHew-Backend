package config

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   *Server
		Auth     *Auth
		Database *Database
	}

	Server struct {
		Port           int
		AllowedOrigins []string
	}
	Auth struct {
		Secret string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		configInstance = &Config{
			Server: &Server{
				Port:           getEnvAsInt("SERVER_PORT"),
				AllowedOrigins: strings.Split(os.Getenv("SERVER_ALLOWED_ORIGINS"), ","),
			},
			Auth: &Auth{
				Secret: os.Getenv("AUTH_SECRET"),
			},
			Database: &Database{
				Host:     os.Getenv("DB_HOST"),
				Port:     getEnvAsInt("DB_PORT"),
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASS"),
				DBName:   os.Getenv("DB_NAME"),
				SSLMode:  os.Getenv("DB_SSLMODE"),
			},
		}
	})
	return configInstance
}

func getEnvAsInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err)
	}
	return value
}
