// Package config provides application configuration logic using Viper.
package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

// Config is the main struct that holds all configuration for the application.
type Config struct {
	Logger   LoggerConfig   `mapstructure:"logger"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// LoggerConfig holds logging-specific settings.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// HTTPConfig holds HTTP server-specific settings.
type HTTPConfig struct {
	Port    string `mapstructure:"port"`
	GinMode string `mapstructure:"gin_mode"`
	BaseURL string `mapstructure:"base_url"`
}

// PostgresConfig holds all settings for the PostgreSQL database connection.
type PostgresConfig struct {
	MasterDSN string     `mapstructure:"master_dsn"`
	Pool      PoolConfig `mapstructure:"pool"`
}

// PoolConfig defines the connection pool settings for the database.
type PoolConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds all settings for the Redis connection.
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NewConfig parses the YAML file and environment variables to return a configuration struct.
func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigFile("configs/config.yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("logger.level", "debug")
	v.SetDefault("http.port", ":8080")
	v.SetDefault("http.gin_mode", "debug")
	v.SetDefault("http.base_url", "http://localhost:8080")
	v.SetDefault("postgres.pool.max_open_conns", 10)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
