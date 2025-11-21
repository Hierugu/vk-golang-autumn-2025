package handlers

import (
	"database/sql"
	"hw5/types"
)

type Handler struct {
	DB     *sql.DB
	Tables map[string][]types.Column
}
