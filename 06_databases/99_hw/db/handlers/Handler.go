package handlers

import (
	"database/sql"
	"hw6/types"
)

type Handler struct {
	DB     *sql.DB
	Tables map[string][]types.Column
}
