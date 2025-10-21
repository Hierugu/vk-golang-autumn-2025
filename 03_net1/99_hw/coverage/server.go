package hw3

import (
	"encoding/json"
	"net/http"
)

const accessKey = "12547890"
const maxLimit = 25

var db Database = XMLDatabase{"dataset.xml"}
var allowedOrderFields = []string{"id", "age", "name", ""}
var allowedOrderByFields = []int{-1, 0, 1}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if !checkHeader(r) {
		http.Error(w, "Wrong access token", http.StatusUnauthorized)
		return
	}

	if !checkMethod(r) {
		http.Error(w, "Use GET method", http.StatusMethodNotAllowed)
		return
	}

	var searchRequest SearchRequest
	var err error
	if searchRequest, err = parseQuery(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchErrorResponse{Error: err.Error()})
		return
	}

	users, err := db.Load()
	if err != nil {
		http.Error(w, "Cannot access database", http.StatusInternalServerError)
		return
	}

	filtered := filterUsers(users, searchRequest.Query)

	if searchRequest.OrderField == "" {
		searchRequest.OrderField = "name"
	}

	sortUsers(filtered, searchRequest.OrderField, searchRequest.OrderBy)

	result := paginateUsers(filtered, searchRequest.Offset, searchRequest.Limit)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
