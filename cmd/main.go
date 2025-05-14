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
		fmt.Println("Usage: go-cli-db A command-line tool to interact with a PostgreSQL database.")
		fmt.Println("\nCommands:")
		fmt.Println("  schemas                List all schemas in the database")
		fmt.Println("  connections            Show active database connections")
		fmt.Println("  size [schema]          List tables sizes in the specified schema")
		fmt.Println("  users                  List database users")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	// Define command-line flags
	homeDir, err := os.UserHomeDir()
	var configPath string
	if err != nil {
		configPath = "config.yaml"
	} else {
		configPath = filepath.Join(homeDir, "config.yaml")
	}

	// Parse the flags
	flag.Parse()

	// Check if the user provided a command
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Get the command from the command line arguments
	command := flag.Arg(0)
	args := flag.Args()[1:]

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		fmt.Println("A config.yaml file is required in the root directory of the project.")
		fmt.Println("The structure of the file should be as follows:")
		fmt.Println("database:")
		fmt.Println("  host: localhost")
		fmt.Println("  port: 5432")
		fmt.Println("  user: postgres")
		fmt.Println("  password: password")
		fmt.Println("  dbname: postgres")
		fmt.Println("  sslmode: disable")
		os.Exit(1)
	}

	// Connect to the database
	database.Connect(cfg.GetDatabaseURL())

	// Handle user commands
	handlers.HandleCommand(command, args)
}
