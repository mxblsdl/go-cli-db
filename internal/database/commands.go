package database

import (
	"fmt"
	"go-cli-db/internal/config"
)

// Add this struct at the top of the file
type Connection struct {
	Username string
	Count    int
}

func GetSchemaNames() error {
	query := `
		select
			distinct
		table_schema
		from
			information_schema.tables
		where
			table_schema not in ('public', 'sde')
		order by
			table_schema;
    `

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return err
	}
	// Print the table
	fmt.Printf("%sScehmas in the database:\n%s", config.Bold, config.Reset)
	fmt.Println("===================================")
	for _, table := range tables {
		fmt.Printf("%sSchema:%s %s\n", config.Green, config.Reset, table)
	}

	return nil
}

func GetActiveConnections() error {
	query := `
		select
			COALESCE(usename, 'none') as username,
			count(*)
		from
			pg_stat_activity
		group by
			usename;
	`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var connectionsTable []Connection
	for rows.Next() {
		var username string
		var count int
		if err := rows.Scan(&username, &count); err != nil {
			return err
		}
		connectionsTable = append(connectionsTable, Connection{Username: username, Count: count})
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// Print the table
	for _, conn := range connectionsTable {
		fmt.Printf("%sUser:%s %s, Count: %d\n", config.Green, config.Reset, conn.Username, conn.Count)
	}

	var totalConnections int
	for _, row := range connectionsTable {
		totalConnections += row.Count

	}
	fmt.Printf("%sActive connections in the database:%s %d\n", config.Red, config.Reset, totalConnections)

	return nil
}

func GetUsers() error {
	query := `
		select
			COALESCE(usename, 'none') as username
		from
			pg_user;
	`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return err
		}
		users = append(users, username)
	}

	if err = rows.Err(); err != nil {
		return err
	}
	// Print the table
	fmt.Printf("%sUsers in the database:%s\n", config.Bold, config.Reset)
	fmt.Println("===================================")
	for _, user := range users {
		fmt.Printf("%sUser:%s %s\n", config.Green, config.Reset, user)
	}

	return nil
}

func GetTableSizes(schema string) error {
	query := fmt.Sprintf(`
	SELECT
		table_name,
		pg_size_pretty(pg_total_relation_size(table_schema || '.' || table_name)) as total_size
	FROM information_schema.tables
	WHERE table_schema = '%s'
	ORDER BY pg_total_relation_size(table_schema || '.' || table_name) DESC;
	`, schema)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		var size string
		if err := rows.Scan(&tableName, &size); err != nil {
			return err
		}
		tables = append(tables, fmt.Sprintf("%s (%s)", tableName, size))
	}

	if err = rows.Err(); err != nil {
		return err
	}
	fmt.Printf("%sTables in the database:%s\n", config.Bold, config.Reset)
	fmt.Println("===================================")
	for _, table := range tables {
		fmt.Printf("%sTable:%s %s\n", config.Green, config.Reset, table)
	}

	return nil
}
