package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   *Server
		Database *Database
		Supabase *Supabase
	}

	Server struct {
		Port           int
		AllowedOrigins []string
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Supabase struct {
		URL string
		Key string
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("[config] no .env file found, using system environment variables")
		}

		configInstance = &Config{
			Server: &Server{
				Port:           getEnvAsInt("SERVER_PORT", 8080),
				AllowedOrigins: strings.Split(getEnv("SERVER_ALLOWED_ORIGINS", "*"), ","),
			},
			Database: &Database{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnvAsInt("DB_PORT", 5432),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASS", ""),
				DBName:   getEnv("DB_NAME", "postgres"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			Supabase: &Supabase{
				URL: getEnv("SUPABASE_URL", ""),
				Key: getEnv("SUPABASE_KEY", ""),
			},
		}
	})
	return configInstance
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
		log.Printf("[config] invalid int for %s=%s, using default %d", key, valueStr, defaultVal)
	}
	return defaultVal
}
