package utils

import (
	"database/sql"
	"fmt"
	"hw5/queries"
	"hw5/types"
)

func LoadTables(db *sql.DB) (map[string][]types.Column, error) {
	result := make(map[string][]types.Column)

	// Read table names first, then close the rows before querying columns for each table.
	table_rows, err := db.Query(queries.GetTables)
	if err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}
	defer table_rows.Close()

	tables := []string{}
	for table_rows.Next() {
		var tn string
		if err := table_rows.Scan(&tn); err != nil {
			return nil, fmt.Errorf("DB error: %w", err)
		}
		tables = append(tables, tn)
	}
	if err := table_rows.Err(); err != nil {
		return nil, fmt.Errorf("DB error: %w", err)
	}

	// Now query columns per table. Outer rows are closed (defer executed at function end,
	// but we don't execute nested queries while iterating an active rows result set).
	for _, table_name := range tables {
		table_cols, err := db.Query(queries.GetTableColumns, table_name)
		if err != nil {
			return nil, fmt.Errorf("DB error: %w", err)
		}

		columns := []types.Column{}
		for table_cols.Next() {
			var col_name, col_type, col_nullable string
			if err := table_cols.Scan(&col_name, &col_type, &col_nullable); err != nil {
				table_cols.Close()
				return nil, fmt.Errorf("DB error: %w", err)
			}
			columns = append(columns, types.Column{Name: col_name, DataType: col_type, Nullable: col_nullable == "YES"})
		}
		if err := table_cols.Err(); err != nil {
			table_cols.Close()
			return nil, fmt.Errorf("DB error: %w", err)
		}
		table_cols.Close()

		result[table_name] = columns
	}

	return result, nil
}
