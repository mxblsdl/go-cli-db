package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Config file not found. %sWould you like to create a new one? (y/n)\n%s", Bold, Reset)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return nil, fmt.Errorf("config file not found")
		}
		// Create a new config file
		err = MakeConfig(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}

		// Try to load the newly created config
		return LoadConfig(path)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode)
}

func MakeConfig(path string) error {
	// Gather database configuration from user
	var host, user, password, dbName, sslMode string
	var port int

	fmt.Print("Database host (default: localhost): ")
	fmt.Scanln(&host)
	if host == "" {
		host = "localhost"
	}

	fmt.Print("Database port (default: 5432): ")
	fmt.Scanln(&port)
	if port == 0 {
		port = 5432
	}

	fmt.Print("Database user: ")
	fmt.Scanln(&user)

	fmt.Print("Database password: ")
	fmt.Scanln(&password)

	fmt.Print("Database name: ")
	fmt.Scanln(&dbName)

	fmt.Print("SSL mode (default: disable): ")
	fmt.Scanln(&sslMode)
	if sslMode == "" {
		sslMode = "disable"
	}

	// Create config content
	configContent := fmt.Sprintf(`database:
  host: "%s"
  port: %d
  user: "%s"
  password: "%s"
  dbname: "%s"
  sslmode: "%s"
`, host, port, user, password, dbName, sslMode)

	// Write config file
	if err := os.WriteFile(path, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("%sConfig file created at%s %s\n", Green, Reset, path)

	return nil
}
