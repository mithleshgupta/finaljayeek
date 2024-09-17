package config

import "os"

type Config struct {
	Host             string
	Port             string
	BaseURL          string
	BaseStorageURL   string
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
	PostgresUsername string
	PostgresPassword string
	PostgresSslMode  string
	PostgresTimeZone string
	RedisHost        string
	RedisPassword    string
	RedisPort        string
	StreamApiKey     string
	StreamApiSecret  string
}

func NewConfig() *Config {
	return &Config{
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		BaseURL:          os.Getenv("BASE_URL"),
		BaseStorageURL:   os.Getenv("BASE_URL") + "/uploads",
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
		PostgresUsername: os.Getenv("POSTGRES_USERNAME"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresSslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		PostgresTimeZone: os.Getenv("POSTGRES_TIME_ZONE"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		StreamApiKey:     os.Getenv("STREAM_API_KEY"),
		StreamApiSecret:  os.Getenv("STREAM_API_SECRET"),
	}
}
