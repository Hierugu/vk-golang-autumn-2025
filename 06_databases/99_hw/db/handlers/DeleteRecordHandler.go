package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) DeleteRecordHandler(w http.ResponseWriter, r *http.Request) {
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

	query := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = ?;", table, pk)
	res, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"response": map[string]interface{}{"deleted": affected}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println("open connections:", h.DB.Stats().OpenConnections)
}
