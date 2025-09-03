package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var configPath string

// SetConfigPath sets the global config path
func SetConfigPath(path string) {
	configPath = path
}

// ListConnections lists all connections in the config file
func ListConnections() error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	fmt.Println("Available database connections:")
	fmt.Println("------------------------------")
	for _, conn := range cfg.Connections {
		fmt.Printf("DefaultConnection: %s\n", cfg.DefaultConnection)
		if conn.Name == cfg.DefaultConnection {
			fmt.Printf("* %s (default) - %s@%s:%d/%s\n",
				conn.Name, conn.User, conn.Host, conn.Port, conn.DBName)
		} else {
			fmt.Printf("  %s - %s@%s:%d/%s\n",
				conn.Name, conn.User, conn.Host, conn.Port, conn.DBName)
		}
	}
	return nil
}

// AddConnection adds a new connection to the config file
func AddConnection() error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		// If config doesn't exist, create a new one
		cfg = &Config{
			Connections: []DatabaseConfig{},
		}
	}

	var name, host, user, password, dbName, sslMode string
	var port int

	fmt.Print("Connection name: ")
	fmt.Scanln(&name)

	// Check if connection name already exists
	for _, conn := range cfg.Connections {
		if conn.Name == name {
			fmt.Println("Connection name already exists. Please choose another name.")
			return nil
		}
	}

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

	// Add the new connection
	newConn := DatabaseConfig{
		Name:     name,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}

	cfg.Connections = append(cfg.Connections, newConn)

	// If this is the first connection, set it as default
	if len(cfg.Connections) == 1 {
		cfg.DefaultConnection = name
	}

	// Save the config
	return SaveConfig(cfg)
}

// SaveConfig saves the config to the file
func SaveConfig(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func EditConnection(name string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	var foundIndex = -1
	for i, conn := range cfg.Connections {
		if conn.Name == name {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("connection not found")
	}

	// Get current connection values to show as defaults
	currentConn := cfg.Connections[foundIndex]

	// Store original name to check against default
	originalName := currentConn.Name
	wasDefault := (cfg.DefaultConnection == originalName)

	var newName, host, user, password, dbName, sslMode string
	var port int

	fmt.Printf("Database connection name (current: %s): ", currentConn.Name)
	fmt.Scanln(&newName)
	if newName == "" {
		newName = currentConn.Name
	}

	if newName != originalName {
		// Check if the new name is already in use
		for _, conn := range cfg.Connections {
			if conn.Name == newName {
				fmt.Printf("Connection name '%s' is already in use\n", newName)
				return nil
			}
		}
	}

	fmt.Printf("Database host (current: %s): ", currentConn.Host)
	fmt.Scanln(&host)
	if host == "" {
		host = currentConn.Host
	}

	fmt.Printf("Database port (current: %d): ", currentConn.Port)
	fmt.Scanln(&port)
	if port == 0 {
		port = currentConn.Port
	}

	fmt.Printf("Database user (current: %s): ", currentConn.User)
	fmt.Scanln(&user)
	if user == "" {
		user = currentConn.User
	}

	fmt.Print("Database password (current: *********): ")
	fmt.Scanln(&password)
	if password == "" {
		password = currentConn.Password
	}

	fmt.Printf("Database name (current: %s): ", currentConn.DBName)
	fmt.Scanln(&dbName)
	if dbName == "" {
		dbName = currentConn.DBName
	}

	fmt.Printf("SSL mode (current: %s): ", currentConn.SSLMode)
	fmt.Scanln(&sslMode)
	if sslMode == "" {
		sslMode = currentConn.SSLMode
	}

	// Update the connection
	cfg.Connections[foundIndex] = DatabaseConfig{
		Name:     newName,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
	if wasDefault {
		cfg.DefaultConnection = newName
		fmt.Printf("Default connection updated to '%s'\n", newName)
	}

	// Save the config
	if err := SaveConfig(cfg); err != nil {
		return err
	}
	fmt.Printf("%sConnection '%s' updated successfully.%s", Green, name, Reset)
	return nil
}

// RemoveConnection removes a connection from the config file
func RemoveConnection(name string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Find the connection to remove
	var foundIndex = -1
	for i, conn := range cfg.Connections {
		if conn.Name == name {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		fmt.Printf("Connection '%s' not found\n", name)
		return nil
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to remove connection '%s'? (y/n): ", name)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Operation cancelled")
		return nil
	}

	// Remove the connection
	cfg.Connections = append(cfg.Connections[:foundIndex], cfg.Connections[foundIndex+1:]...)

	// If we removed the default connection, set a new default if available
	if cfg.DefaultConnection == name && len(cfg.Connections) > 0 {
		cfg.DefaultConnection = cfg.Connections[0].Name
		fmt.Printf("Default connection set to '%s'\n", cfg.DefaultConnection)
	}

	// Save the config
	if err := SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("%sConnection '%s' removed successfully%s\n", Green, name, Reset)
	return nil
}

// SetDefaultConnection sets the default connection in the config
func SetDefaultConnection(name string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Check if the connection exists
	var found bool
	for _, conn := range cfg.Connections {
		if conn.Name == name {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Connection '%s' not found\n", name)
		return nil
	}

	// Set as default
	cfg.DefaultConnection = name

	// Save the config
	if err := SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("%sDefault connection set to '%s'%s\n", Green, name, Reset)
	return nil
}
