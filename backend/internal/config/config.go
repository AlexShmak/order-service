package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"slices"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"local"`
	Database    DatabaseConfig
	Server      ServerConfig
	JWT         JWT
	Kafka       KafkaConfig
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" env-required:"true"`
	Topic   string   `env:"KAFKA_TOPIC" env-required:"true"`
}

type JWT struct {
	AccessSecret  string `env:"ACCESS_SECRET" env-required:"true"`
	RefreshSecret string `env:"REFRESH_SECRET" env-required:"true"`
}

type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" env-default:"localhost"`
	Port            string        `env:"DB_PORT" env-default:"5432"`
	User            string        `env:"DB_USER" env-required:"true"`
	Password        string        `env:"DB_PASSWORD" env-required:"true"`
	DBName          string        `env:"DB_NAME" env-required:"true"`
	SSLMode         string        `env:"SSL_MODE" env-default:"disable"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" env-default:"25"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"5m"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" env-default:"1m"`
	URL             string        `env:"PGSQL_URL" env-required:"true"`
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT" env-default:"8080"`
	Host         string        `env:"SERVER_HOST" env-default:"localhost"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"30s"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &cfg, nil
}

// Basic custom validation for config
func (c *Config) validate() error {
	validEnvironments := []string{"local", "dev", "prod"}

	if !slices.Contains(validEnvironments, c.Environment) {
		return fmt.Errorf("invalid environment: %s, must be one of %v", c.Environment, validEnvironments)
	}

	validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
	if !slices.Contains(validSSLModes, c.Database.SSLMode) {
		return fmt.Errorf("invalid SSL mode: %s, must be one of %v", c.Database.SSLMode, validSSLModes)
	}

	if c.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("max open connections must be positive, got: %d", c.Database.MaxOpenConns)
	}

	if c.Database.MaxIdleConns <= 0 {
		return fmt.Errorf("max idle connections must be positive, got: %d", c.Database.MaxIdleConns)
	}

	return nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
