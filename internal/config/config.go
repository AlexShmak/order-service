package config

import (
	"fmt"
	"slices"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string         `env:"ENV" env-default:"local"`
	Database    DatabaseConfig `env-prefix:"DB_"`
	Server      ServerConfig   `env-prefix:"SERVER_"`
}

type DatabaseConfig struct {
	Host            string        `env:"HOST" env-default:"localhost"`
	Port            string        `env:"PORT" env-default:"5432"`
	DBName          string        `env:"DBNAME" env-required:"true"`
	AdminUser       string        `env:"ADMIN_USER" env-required:"true"`
	AdminPassword   string        `env:"ADMIN_PASSWORD" env-required:"true"`
	RegularUser     string        `env:"REGULAR_USER" env-required:"true"`
	RegularPassword string        `env:"REGULAR_PASSWORD" env-required:"true"`
	SSLMode         string        `env:"SSL_MODE" env-default:"disable"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" env-default:"25"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"5m"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" env-default:"1m"`
}

type ServerConfig struct {
	Port         string        `env:"PORT" env-default:"8080"`
	Host         string        `env:"HOST" env-default:"localhost"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"30s"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
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

func (c *Config) GetAdminPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.AdminUser,
		c.Database.AdminPassword,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRegularPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.RegularUser,
		c.Database.RegularPassword,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
