package handlers

import (
	"fmt"
	"go-cli-db/internal/config"
	"go-cli-db/internal/database"
)

// HandleCommand processes the user command and interacts with the database.
func HandleCommand(command string, args []string) {
	switch command {
	case "connections":
		// Call the function to list items from the database
		err := database.GetActiveConnections()
		if err != nil {
			fmt.Printf("%sError getting active connections:%s %s", config.Red, config.Reset, err)
			return
		}
	case "schemas":
		err := database.GetSchemaNames()
		if err != nil {
			fmt.Printf("%sError getting schema names:%s %s", config.Red, config.Reset, err)
			return
		}
	case "users":
		err := database.GetUsers()
		if err != nil {
			fmt.Printf("%sError getting users:%s %s", config.Red, config.Reset, err)
			return
		}
	case "size":
		if len(args) < 1 {
			// fmt.Printf("%sError:%s missing schema name argument for 'size' command", config.Red, config.Reset)
			// return
			database.GetSchemaSizes()
			return
		}
		err := database.GetTableSizes(args[0])
		if err != nil {
			fmt.Printf("%sError getting schema size:%s %s", config.Red, config.Reset, err)
			return
		}

	case "config":
		if len(args) == 0 {
			fmt.Println("Available config commands:")
			fmt.Println("  config list         - List available connections")
			fmt.Println("  config add          - Add a new connection")
			fmt.Println("  config edit [name]  - Edit a connection")
			fmt.Println("  config remove [name]- Remove a connection")
			fmt.Println("  config use [name]   - Set default connection")
			return
		}

		configCmd := args[0]
		configArgs := args[1:]

		switch configCmd {
		case "list":
			config.ListConnections()
		case "add":
			config.AddConnection()
		case "edit":
			if len(configArgs) == 0 {
				fmt.Println("Please specify a connection name to edit")
				return
			}
			config.EditConnection(configArgs[0])
		case "remove":
			if len(configArgs) == 0 {
				fmt.Println("Please specify a connection name to remove")
				return
			}
			config.RemoveConnection(configArgs[0])
		case "use":
			if len(configArgs) == 0 {
				fmt.Println("Please specify a connection name to set as default")
				return
			}
			config.SetDefaultConnection(configArgs[0])
		default:
			fmt.Printf("Unknown config command: %s\n", configCmd)
		}
	default:
		fmt.Printf("%sUnknown command:%s %s", config.Red, config.Reset, command)
		return
	}
}
