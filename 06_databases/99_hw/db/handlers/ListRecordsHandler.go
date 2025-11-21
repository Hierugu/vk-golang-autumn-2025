package handlers

import (
	"encoding/json"
	"fmt"
	"hw5/queries"
	"net/http"
	"strconv"
)

func (h *Handler) ListRecordsHandler(w http.ResponseWriter, r *http.Request) {
	limit, offset := 5, 0

	table := r.PathValue("table")
	if h.Tables[table] == nil {
		http.Error(w, `{"error": "unknown table"}`, http.StatusNotFound)
		return
	}

	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v >= 0 {
		limit = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v >= 0 {
		offset = v
	}
	fmt.Println(table, limit, offset)

	rows, err := h.DB.Query(fmt.Sprintf(queries.ListRecords, table), limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	records := []interface{}{}
	cols, err := rows.Columns()
	if err != nil {
		http.Error(w, `{"error": "failed to read columns metadata"}`, http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			http.Error(w, `{"error": "failed to scan record from DB"}`, http.StatusInternalServerError)
			return
		}

		rec := make(map[string]interface{})
		for i, v := range vals {
			switch tv := v.(type) {
			case []byte:
				rec[cols[i]] = string(tv)
			default:
				rec[cols[i]] = tv
			}
		}

		records = append(records, rec)
	}

	resp := map[string]interface{}{"response": map[string]interface{}{"records": records}}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	respJSON, _ := json.Marshal(resp)
	w.Write(respJSON)
}
