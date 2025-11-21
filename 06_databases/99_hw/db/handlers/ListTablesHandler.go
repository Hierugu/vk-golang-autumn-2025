package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
)

func (h *Handler) ListTablesHandler(w http.ResponseWriter, r *http.Request) {
	tables := []string{}
	for tableName := range h.Tables {
		tables = append(tables, tableName)
	}
	sort.Strings(tables)

	resp := map[string]interface{}{"response": map[string]interface{}{"tables": tables}}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(resp)
	w.Write(respJSON)

	// debug print removed
}
