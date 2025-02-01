package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	User string
	Password string
	Host string
	Port string
	Name string
}

type Config struct {
	Database DatabaseConfig
}

func Load() (*Config, error) {
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		Database: *dbConfig,
	}, nil
}

func loadDatabaseConfig() (*DatabaseConfig, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}




	config := &DatabaseConfig{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
	}

	err = validateDatabaseConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateDatabaseConfig(config *DatabaseConfig) error {

	if config.User == "" {
		return fmt.Errorf("database user is required")
	}

	if config.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if config.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Port == "" {
		return fmt.Errorf("database port is required")
	}

	if config.Name == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
	
}



