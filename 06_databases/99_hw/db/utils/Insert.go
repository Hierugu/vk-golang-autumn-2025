package utils

import (
	"database/sql"
	"fmt"
	"hw5/queries"
	"sort"
)

func Insert(db *sql.DB, tableName string, values map[string]interface{}) (int, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("no values to insert")
	}

	// Collect and sort column names for deterministic order
	cols := make([]string, 0, len(values))
	for k := range values {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	// Build column list and placeholders, and args slice
	colList := ""
	placeholders := ""
	args := make([]interface{}, 0, len(cols))
	for i, c := range cols {
		if i > 0 {
			colList += ", "
			placeholders += ", "
		}
		// quote identifiers with backticks to be safe for MySQL
		colList += fmt.Sprintf("`%s`", c)
		placeholders += "?"
		args = append(args, values[c])
	}

	// Quote table name as well
	table := fmt.Sprintf("`%s`", tableName)
	query := fmt.Sprintf(queries.InsertRecord, table, colList, placeholders)

	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("insert error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}
