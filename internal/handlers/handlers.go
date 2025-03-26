package handlers

import (
	"fmt"
	"go-cli-db/internal/database"
	"os"
)

// HandleCommand processes the user command and interacts with the database.
func HandleCommand(command string) {
	switch command {
	case "connections":
		// Call the function to list items from the database
		cons, err := database.GetActiveConnections()
		if err != nil {
			fmt.Println("Error getting active connections:", err)
			os.Exit(1)
		}
		fmt.Println("Active connections in the database:", cons)
	case "schemas":
		tables, err := database.GetSchemaNames()
		if err != nil {
			fmt.Println("Error getting schema names:", err)
			os.Exit(1)
		}
		fmt.Println("Scehmas in the database:")
		for _, table := range tables {
			fmt.Println(table)
		}
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
