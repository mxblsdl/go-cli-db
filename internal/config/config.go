package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	DefaultConnection string           `yaml:"default_connection"`
	Connections       []DatabaseConfig `yaml:"connections"`
}

func (c *Config) GetDatabaseURL(connectionName string) string {
	if connectionName == "" {
		connectionName = c.DefaultConnection
	}

	for _, conn := range c.Connections {
		if conn.Name == connectionName {
			return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				conn.Host, conn.Port, conn.User, conn.Password, conn.DBName, conn.SSLMode)
		}
	}

	return ""
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{
		Connections: []DatabaseConfig{},
	}

	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	if len(file) == 0 {
		return config, nil
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	if len(config.Connections) > 0 {
		defaultExists := false

		for _, conn := range config.Connections {
			if conn.Name == config.DefaultConnection {
				defaultExists = true
				break
			}
		}

		// Ensure there is a default connection
		if !defaultExists || config.DefaultConnection == "" {
			config.DefaultConnection = config.Connections[0].Name
		}
	}

	return config, nil
}
