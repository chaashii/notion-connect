package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	// Database DatabaseConfig
	// JWT      JWTConfig
	// Redis    RedisConfig
	NotionAPI NotionAPIConfig
	Env       string
}

type ServerConfig struct {
	Port    string
	Host    string
	Timeout time.Duration
	Mode    string
}

// type DatabaseConfig struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	DBName   string
// 	SSLMode  string
// }

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// type RedisConfig struct {
// 	Host     string
// 	Port     string
// 	Password string
// 	DB       int
// }

type NotionAPIConfig struct {
	Key  string
	DbId string
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Set default configurations
	setDefaults()

	// Read .env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %s", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port:    viper.GetString("SERVER_PORT"),
			Host:    viper.GetString("SERVER_HOST"),
			Timeout: viper.GetDuration("SERVER_TIMEOUT") * time.Second,
			Mode:    viper.GetString("GIN_MODE"),
		},
		// Database: DatabaseConfig{
		// 	Host:     viper.GetString("DB_HOST"),
		// 	Port:     viper.GetString("DB_PORT"),
		// 	User:     viper.GetString("DB_USER"),
		// 	Password: viper.GetString("DB_PASSWORD"),
		// 	DBName:   viper.GetString("DB_NAME"),
		// 	SSLMode:  viper.GetString("DB_SSL_MODE"),
		// },
		// JWT: JWTConfig{
		// 	Secret:          viper.GetString("JWT_SECRET"),
		// 	ExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
		// },
		// Redis: RedisConfig{
		// 	Host:     viper.GetString("REDIS_HOST"),
		// 	Port:     viper.GetString("REDIS_PORT"),
		// 	Password: viper.GetString("REDIS_PASSWORD"),
		// 	DB:       viper.GetInt("REDIS_DB"),
		// },
		NotionAPI: NotionAPIConfig{
			Key:  viper.GetString("NOTION_API_KEY"),
			DbId: viper.GetString("NOTION_DATABASE_ID"),
		},
		Env: viper.GetString("ENV"),
	}

	// Validate required configurations
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_TIMEOUT", 30)
	viper.SetDefault("GIN_MODE", "debug")

	// Database defaults
	// viper.SetDefault("DB_PORT", "5432")
	// viper.SetDefault("DB_SSL_MODE", "disable")

	// JWT defaults
	// viper.SetDefault("JWT_EXPIRATION_HOURS", 24)

	// Redis defaults
	// viper.SetDefault("REDIS_PORT", "6379")
	// viper.SetDefault("REDIS_DB", 0)

	// Environment default
	viper.SetDefault("ENV", "development")
}

func validateConfig(config *Config) error {
	// 	// Add validation for required fields
	// 	if config.Database.Password == "" {
	// 		return fmt.Errorf("database password is required")
	// 	}

	// 	if config.JWT.Secret == "" {
	// 		return fmt.Errorf("JWT secret is required")
	// 	}

	if config.NotionAPI.DbId == "" || config.NotionAPI.Key == "" {
		return fmt.Errorf("API key is required")
	}

	return nil
}

// GetDSN returns the database connection string
// func (c *DatabaseConfig) GetDSN() string {
// 	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
// }

// IsDevelopment checks if the current environment is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}
