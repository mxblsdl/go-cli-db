package handlers

import (
	"fmt"
	"go-cli-db/internal/config"
	"go-cli-db/internal/database"
	"os"
)

// HandleCommand processes the user command and interacts with the database.
func HandleCommand(command string, args []string) {
	switch command {
	case "connections":
		// Call the function to list items from the database
		err := database.GetActiveConnections()
		if err != nil {
			fmt.Printf("%sError getting active connections:%s %s", config.Red, config.Reset, err)
			os.Exit(1)
		}
	case "schemas":
		err := database.GetSchemaNames()
		if err != nil {
			fmt.Printf("%sError getting schema names:%s %s", config.Red, config.Reset, err)
			os.Exit(1)
		}
	case "users":
		err := database.GetUsers()
		if err != nil {
			fmt.Printf("%sError getting users:%s %s", config.Red, config.Reset, err)
			os.Exit(1)
		}
	case "size":
		if len(args) < 1 {
			fmt.Printf("%sError:%s missing schema name argument for 'size' command", config.Red, config.Reset)
			os.Exit(1)
		}
		err := database.GetTableSizes(args[0])
		if err != nil {
			fmt.Printf("%sError getting schema size:%s %s", config.Red, config.Reset, err)
			os.Exit(1)
		}
	default:
		fmt.Printf("%sUnknown command:%s %s", config.Red, config.Reset, command)
		os.Exit(1)
	}
}
