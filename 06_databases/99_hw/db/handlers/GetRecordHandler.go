package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) GetRecordHandler(w http.ResponseWriter, r *http.Request) {
	table, id := r.PathValue("table"), r.PathValue("id")

	if h.Tables[table] == nil {
		http.Error(w, `{"error": "unknown table"}`, http.StatusNotFound)
		return
	}

	// determine primary key column: prefer "id", then any column containing "id"
	pk := "id"
	found := false
	for _, c := range h.Tables[table] {
		if c.Name == "id" {
			pk = "id"
			found = true
			break
		}
	}
	if !found {
		for _, c := range h.Tables[table] {
			if strings.HasSuffix(c.Name, "_id") || strings.Contains(c.Name, "id") {
				pk = c.Name
				found = true
				break
			}
		}
	}

	// Build query: SELECT * FROM `table` WHERE `pk` = ? LIMIT 1
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s` = ? LIMIT 1;", table, pk)

	// Use Query to get Columns metadata and row data.
	rows, err := h.DB.Query(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, `{"error": "record not found"}`, http.StatusNotFound)
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		http.Error(w, `{"error": "failed to read columns metadata"}`, http.StatusInternalServerError)
		return
	}

	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	if err := rows.Scan(ptrs...); err != nil && err != sql.ErrNoRows {
		http.Error(w, `{"error": "failed to scan record from DB"}`, http.StatusInternalServerError)
		return
	}

	record := make(map[string]interface{})
	for i, c := range cols {
		v := vals[i]
		switch tv := v.(type) {
		case []byte:
			record[c] = string(tv)
		default:
			record[c] = tv
		}
	}

	response := map[string]interface{}{"response": map[string]interface{}{"record": record}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println("open connections:", h.DB.Stats().OpenConnections)
}
