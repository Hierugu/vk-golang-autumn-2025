package handlers

import (
	"encoding/json"
	"fmt"
	"hw6/utils"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) CreateRecordHandler(w http.ResponseWriter, r *http.Request) {
	table := r.PathValue("table")
	if h.Tables[table] == nil {
		http.Error(w, `{"error": "unknown table"}`, http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read request body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var f interface{}
	err = json.Unmarshal(body, &f)

	if err != nil {
		http.Error(w, `{"error": "invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	record, ok := f.(map[string]interface{})
	if !ok {
		http.Error(w, `{"error": "expected JSON object"}`, http.StatusBadRequest)
		return
	}
	recordValidated, err := utils.Validate(record, h.Tables[table])
	if err != nil {
		fmt.Println("validation error:", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// If the table has a primary column named exactly "id", ignore client-supplied value
	// so AUTO_INCREMENT will assign the value (test expects this behavior for `items` table).
	for _, c := range h.Tables[table] {
		if c.Name == "id" {
			delete(recordValidated, "id")
			break
		}
	}

	id, err := utils.Insert(h.DB, table, recordValidated)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// determine primary key column name (prefer exact "id", otherwise use first *_id or containing "id")
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

	// If client provided pk value in the validated record, use it; otherwise use last insert id
	var respVal interface{} = id
	if v, ok := recordValidated[pk]; ok {
		respVal = v
	}

	resp := map[string]interface{}{"response": map[string]interface{}{pk: respVal}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// debug print removed
}
