package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
		SSLMode  string
	}

	JWT struct {
		Secret             string
		TokenExpiry        time.Duration
		RefreshTokenExpiry time.Duration
	}

	Environment string
}

func Load() (*Config, error) {
	godotenv.Load() // Load .env if exists

	cfg := &Config{}

	// Server config
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	// Database config
	cfg.Database.Host = getEnv("DATABASE_HOST", "localhost")
	cfg.Database.Port = getEnv("DATABASE_PORT", "5432")
	cfg.Database.Username = getEnv("DATABASE_USERNAME", "novita")
	cfg.Database.Password = getEnv("DB_PASSWORD", "root")
	cfg.Database.DBName = getEnv("DB_NAME", "shelter_it")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT config
	cfg.JWT.Secret = getEnv("JWT_SECRET", "shelter-it-secret-key")
	cfg.JWT.TokenExpiry = time.Hour * 24         // 24h
	cfg.JWT.RefreshTokenExpiry = time.Hour * 168 // 7d

	cfg.Environment = getEnv("ENVIRONMENT", "dev")

	return cfg, nil
}

func getEnv(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode)
}
