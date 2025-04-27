package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Dynamodb struct {
	Region       string
	Endpoint     string
	DaxHostPorts []string
	UseDax       bool
}
type Config struct {
	AppEnv         string
	DatabaseURL    string
	ServerPort     string
	CacheTableName string
	Dynamodb       Dynamodb
}

func LoadConfig() (*Config, error) {
	useDax, _ := strconv.ParseBool("DYNAMODB_USE_DAX")
	cfg := &Config{
		AppEnv:         getEnv("ENV", ""),
		DatabaseURL:    getEnv("DSN", ""),
		ServerPort:     getEnv("PORT", ""),
		CacheTableName: getEnv("DYNAMODB_CACHE_TABLE_NAME", ""),
		Dynamodb: Dynamodb{
			Region:       getEnv("DYNAMODB_REGION", ""),
			Endpoint:     getEnv("DYNAMODB_ENDPOINT", ""),
			UseDax:       useDax,
			DaxHostPorts: getEnvAsSlice("DYNAMODB_DAX_HOST_PORTS", []string{}),
		},
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("missing required environment variable: DATABASE_URL")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		return defaultVal
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i]) // מנקה רווחים מיותרים
	}
	return parts
}
