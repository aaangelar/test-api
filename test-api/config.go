package dataexport

import "os"

// Config data struct
type Config struct {
	APIPort          string
	ConnectionString string
	DBType           string
}

// NewConfig returns an object of Config data struct with default values from environment variables
func NewConfig() *Config {

	cfg := Config{
		APIPort:          os.Getenv("API_PORT"),
		ConnectionString: os.Getenv("DB_CONNECTION"),
		DBType:           os.Getenv("DB_TYPE"),
	}

	return &cfg
}
