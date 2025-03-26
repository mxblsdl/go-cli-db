package database

import (
	"fmt"
)

// Add this struct at the top of the file
type Connection struct {
	Username string
	Count    int
}

func GetSchemaNames() ([]string, error) {
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
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func GetActiveConnections() (int, error) {
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
		return 0, err
	}
	defer rows.Close()

	var connectionsTable []Connection
	for rows.Next() {
		var username string
		var count int
		if err := rows.Scan(&username, &count); err != nil {
			return 0, err
		}
		connectionsTable = append(connectionsTable, Connection{Username: username, Count: count})
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	// Print the table (optional, for debugging)
	for _, conn := range connectionsTable {
		// fmt.Println(row)
		fmt.Printf("User: %s, Count: %d\n", conn.Username, conn.Count)
	}

	var totalConnections int
	for _, row := range connectionsTable {
		totalConnections += row.Count

	}

	return totalConnections, nil
}
