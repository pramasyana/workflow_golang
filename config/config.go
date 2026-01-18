package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
}

// DatabaseConfig holds database configuration (MySQL)
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parse_time"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"`
	Issuer     string `yaml:"issuer"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DSN returns the MySQL connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=UTC",
		d.Username, d.Password, d.Host, d.Port, d.Name, d.Charset, d.ParseTime)
}

// GetConnMaxLifetime returns connection max lifetime as time.Duration
func (d *DatabaseConfig) GetConnMaxLifetime() time.Duration {
	return time.Duration(d.ConnMaxLifetime) * time.Second
}

// Load loads configuration from the specified file
// It first loads environment variables from .env file if it exists
// then loads configuration from YAML file, and finally overrides with environment variables
func Load(path string) (*Config, error) {
	// Try to load .env file (silent if not found)
	_ = godotenv.Load()

	// If config file exists, load from it
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		// Override with environment variables
		cfg.overrideFromEnv()
		return &cfg, nil
	}

	// If no config file, create config from environment variables only
	cfg := &Config{
		App: AppConfig{
			Host: getEnvString("APP_HOST", "0.0.0.0"),
			Port: getEnvInt("APP_PORT", 8080),
			Name: getEnvString("APP_NAME", "Workflow Approval System"),
			Env:  getEnvString("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnvString("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 3306),
			Username:        getEnvString("DB_USERNAME", "root"),
			Password:        getEnvString("DB_PASSWORD", "password"),
			Name:            getEnvString("DB_NAME", "workflow_approval"),
			Charset:         getEnvString("DB_CHARSET", "utf8mb4"),
			ParseTime:       getEnvBool("DB_PARSE_TIME", true),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 300),
		},
		JWT: JWTConfig{
			Secret:     getEnvString("JWT_SECRET", "your-super-secret-key-change-in-production"),
			Expiration: getEnvInt("JWT_EXPIRATION", 24),
			Issuer:     getEnvString("JWT_ISSUER", "workflow-approval-system"),
		},
		Logging: LoggingConfig{
			Level:  getEnvString("LOGGING_LEVEL", "debug"),
			Format: getEnvString("LOGGING_FORMAT", "json"),
		},
	}

	return cfg, nil
}

// overrideFromEnv allows environment variables to override config values
func (c *Config) overrideFromEnv() {
	// App config
	if host := os.Getenv("APP_HOST"); host != "" {
		c.App.Host = host
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &c.App.Port)
	}
	if env := os.Getenv("APP_ENV"); env != "" {
		c.App.Env = env
	}

	// Database config
	if host := os.Getenv("DB_HOST"); host != "" {
		c.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &c.Database.Port)
	}
	if user := os.Getenv("DB_USERNAME"); user != "" {
		c.Database.Username = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		c.Database.Password = pass
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		c.Database.Name = name
	}
	if charset := os.Getenv("DB_CHARSET"); charset != "" {
		c.Database.Charset = charset
	}
	if parseTime := os.Getenv("DB_PARSE_TIME"); parseTime != "" {
		c.Database.ParseTime = parseTime == "true" || parseTime == "1"
	}
	if maxOpen := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpen != "" {
		fmt.Sscanf(maxOpen, "%d", &c.Database.MaxOpenConns)
	}
	if maxIdle := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdle != "" {
		fmt.Sscanf(maxIdle, "%d", &c.Database.MaxIdleConns)
	}
	if connLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); connLifetime != "" {
		fmt.Sscanf(connLifetime, "%d", &c.Database.ConnMaxLifetime)
	}

	// JWT config
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		c.JWT.Secret = secret
	}
	if exp := os.Getenv("JWT_EXPIRATION"); exp != "" {
		fmt.Sscanf(exp, "%d", &c.JWT.Expiration)
	}
	if issuer := os.Getenv("JWT_ISSUER"); issuer != "" {
		c.JWT.Issuer = issuer
	}

	// Logging config
	if level := os.Getenv("LOGGING_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	if format := os.Getenv("LOGGING_FORMAT"); format != "" {
		c.Logging.Format = format
	}
}

// getEnvString returns environment variable or default value
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as int or default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool returns environment variable as bool or default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

// Address returns the server address
func (a *AppConfig) Address() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
