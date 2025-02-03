package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
	_ "github.com/lib/pq"

)

type DatabaseConfig struct {
	User string
	Password string
	Host string
	Port string
	Name string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnectionMaxLifetime time.Duration
}

type RedisConfig struct {
	Host string
	Port string
	Password string
	DB int
}

type JWTConfig struct {
	Secret string
	RefreshSecret string
	AccessTTL time.Duration
	RefreshTTL time.Duration
}

type Config struct {
	Database *DatabaseConfig
	Redis *RedisConfig
	JWT *JWTConfig


}

func Load() (*Config, error) {

	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	redisConfig, err := loadRedisConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := loadJWTConfig()
	if err != nil {
		return nil, err
	}


	return &Config{
		Database: dbConfig,
		Redis: redisConfig,
		JWT: jwtConfig,
	}, nil

}

func loadDatabaseConfig() (*DatabaseConfig, error) {

	config := &DatabaseConfig{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
		MaxOpenConnections: 25,
		MaxIdleConnections: 5,
		ConnectionMaxLifetime: 5 * time.Minute,
	}

	err := validateDatabaseConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadJWTConfig() (*JWTConfig, error){
	config :=&JWTConfig{
		Secret: os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AccessTTL: 15 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
	}

	err := config.validateJWTConfig()
	if err != nil {
		return nil, err
	}
	return config, nil


}

func loadRedisConfig() (*RedisConfig, error) {

	config := &RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	}

	err := config.validateRedisConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (d *DatabaseConfig) GetDSN() string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
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

	if _, err := strconv.Atoi(config.Port); err != nil {
		return fmt.Errorf("database port is required")
	}

	if config.Name == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
	
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r RedisConfig) validateRedisConfig() error {

	if r.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	if _, err := strconv.Atoi(r.Port); err != nil {
		return fmt.Errorf("redis port is required")
	}
	return nil
}
func (j JWTConfig) validateJWTConfig() error {

	if j.Secret == "" {
		return fmt.Errorf("secret key is required")
	}

	if j.RefreshSecret == "" {
		return fmt.Errorf("refresh secret key is required")
	}

	return nil
}



