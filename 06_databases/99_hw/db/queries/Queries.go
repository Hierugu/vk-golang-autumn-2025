package queries

const GetTables = "SHOW TABLES;"
const GetTableColumns = `
	SELECT
		COLUMN_NAME AS col_name,
		COLUMN_TYPE AS col_type,
		IS_NULLABLE AS col_nullable
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_SCHEMA = DATABASE()
	AND TABLE_NAME = ?;`
const ListRecords = "SELECT * FROM %s LIMIT ? OFFSET ?;"
const InsertRecord = "INSERT INTO %s (%s) VALUES (%s);"
