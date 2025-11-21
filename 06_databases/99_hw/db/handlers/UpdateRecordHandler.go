package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) UpdateRecordHandler(w http.ResponseWriter, r *http.Request) {
	table, id := r.PathValue("table"), r.PathValue("id")

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
	if err := json.Unmarshal(body, &f); err != nil {
		http.Error(w, `{"error": "invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	updates, ok := f.(map[string]interface{})
	if !ok {
		http.Error(w, `{"error": "expected JSON object"}`, http.StatusBadRequest)
		return
	}

	// determine primary key name
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

	// cannot update primary key
	if _, ok := updates[pk]; ok {
		http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, pk), http.StatusBadRequest)
		return
	}

	// build SET clause and args
	setParts := []string{}
	args := []interface{}{}

	for k, v := range updates {
		// ignore unknown columns
		var foundCol *struct {
			Name     string
			DataType string
			Nullable bool
		}
		for _, c := range h.Tables[table] {
			if c.Name == k {
				foundCol = &struct {
					Name     string
					DataType string
					Nullable bool
				}{Name: c.Name, DataType: c.DataType, Nullable: c.Nullable}
				break
			}
		}
		if foundCol == nil {
			// ignore unknown fields
			continue
		}

		// validate value type similar to utils.Validate
		dt := strings.ToLower(foundCol.DataType)
		var val interface{}
		if v == nil {
			if !foundCol.Nullable {
				http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
				return
			}
			val = nil
		} else if strings.HasPrefix(dt, "varchar") || dt == "text" {
			// For string-like columns accept only string or []byte.
			// If a number (JSON float64) is provided — return 400 (do not coerce to string).
			switch t := v.(type) {
			case string:
				val = t
			case []byte:
				val = string(t)
			default:
				// reject numeric and other types for string columns
				http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
				return
			}
		} else if dt == "int" || strings.HasPrefix(dt, "int") {
			switch t := v.(type) {
			case float64:
				if t == float64(int64(t)) {
					val = int(t)
				} else {
					http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
					return
				}
			case string:
				iv, err := strconv.Atoi(t)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
					return
				}
				val = iv
			case []byte:
				iv, err := strconv.Atoi(string(t))
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
					return
				}
				val = iv
			default:
				http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
				return
			}
		} else if dt == "float" || strings.HasPrefix(dt, "float") || strings.HasPrefix(dt, "double") {
			switch t := v.(type) {
			case float64:
				val = t
			case string:
				fv, err := strconv.ParseFloat(t, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
					return
				}
				val = fv
			case []byte:
				fv, err := strconv.ParseFloat(string(t), 64)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
					return
				}
				val = fv
			default:
				http.Error(w, fmt.Sprintf(`{"error": "field %s have invalid type"}`, k), http.StatusBadRequest)
				return
			}
		} else {
			// fallback: accept value as-is
			val = v
		}

		setParts = append(setParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, val)
	}

	if len(setParts) == 0 {
		// nothing to update
		http.Error(w, `{"error": "no valid fields to update"}`, http.StatusBadRequest)
		return
	}

	// append pk value as last arg
	args = append(args, id)

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE `%s` = ?;", table, strings.Join(setParts, ", "), pk)
	res, err := h.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"response": map[string]interface{}{"updated": affected}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println("open connections:", h.DB.Stats().OpenConnections)
}
