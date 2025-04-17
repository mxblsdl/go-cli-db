package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-cli-db/internal/config"
	"go-cli-db/internal/database"
	"go-cli-db/internal/handlers"
)

func main() {
	flag.Usage = func() {
		log.Println("Usage: go-cli-db A command-line tool to interact with a PostgreSQL database.")
		flag.PrintDefaults()
	}

	// Define command-line flags
	configPath := flag.String("config", "config.yaml", "../config.yaml")

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
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
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
