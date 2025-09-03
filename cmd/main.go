package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go-cli-db/internal/config"
	"go-cli-db/internal/database"
	"go-cli-db/internal/handlers"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: go-cli-db [options] command [arguments]")
		fmt.Println("\nA command-line tool to interact with a PostgreSQL database.")
		fmt.Println("\nCommands:")
		fmt.Println("  schemas                List all schemas in the database")
		fmt.Println("  connections            Show active database connections")
		fmt.Println("  size [schema]          List tables sizes in the specified schema")
		fmt.Println("  users                  List database users")
		fmt.Println("  config [subcommand]    Manage database connections")
		fmt.Println("\nOptions:")
		fmt.Println("  -connection string     Name of the database connection to use (overrides default)")
		fmt.Println("                         Run 'db config list' to see available connections")
		// flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  db schemas                        # List schemas using default connection")
		fmt.Println("  db -connection=prod users         # List users on the 'prod' connection")
		fmt.Println("  db config list                    # List all configured connections")
	}

	// Define command-line flags
	homeDir, err := os.UserHomeDir()
	var configPath string
	if err != nil {
		configPath = "dbconfig.yaml"
	} else {
		configPath = filepath.Join(homeDir, "dbconfig.yaml")
	}
	// fmt.Printf("Config path: %s\n", configPath)

	// Parse the flags
	connectionName := flag.String("connection", "", "Name of the database connection to use (overrides default in config)")
	flag.Parse()

	// Check if the user provided a command
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}
	// Set global config path
	config.SetConfigPath(configPath)

	// Get the command from the command line arguments
	command := flag.Arg(0)
	args := flag.Args()[1:]

	// special handling for the config command
	if command == "config" {
		handlers.HandleCommand(command, args)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		fmt.Println("Run 'db config add' to create a new config file.")

		os.Exit(0)
	}

	connString := cfg.GetDatabaseURL(*connectionName)
	if connString == "" {
		fmt.Println("Connection not found")
		os.Exit(0)
	}

	// Connect to the database
	database.Connect(connString)

	// Handle user commands
	handlers.HandleCommand(command, args)
}
