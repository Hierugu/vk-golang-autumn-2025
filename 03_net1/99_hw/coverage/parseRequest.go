package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
)

func checkHeader(r *http.Request) bool {
	return r.Header.Get("AccessToken") == accessKey
}

func checkMethod(r *http.Request) bool {
	return r.Method == http.MethodGet
}

func parseQuery(r *http.Request) (SearchRequest, error) {
	sr := SearchRequest{}
	query := r.URL.Query()
	var err error

	// Query
	sr.Query = query.Get("query")

	// OrderField
	sr.OrderField = query.Get("order_field")
	if !slices.Contains(allowedOrderFields, sr.OrderField) {
		return SearchRequest{}, fmt.Errorf(ErrorBadOrderField)
	}

	// OrderBy
	if s := query.Get("order_by"); s != "" {
		if sr.OrderBy, err = strconv.Atoi(s); err != nil {
			return SearchRequest{}, fmt.Errorf("Invalid order_by: %v", err)
		}
		if !slices.Contains(allowedOrderByFields, sr.OrderBy) {
			return SearchRequest{}, fmt.Errorf("Invalid order_by: must be one of -1, 0, 1")
		}
	} else {
		sr.OrderBy = OrderByAsIs
	}

	// Offset
	if s := query.Get("offset"); s != "" {
		if sr.Offset, err = strconv.Atoi(s); err != nil {
			return SearchRequest{}, fmt.Errorf("Invalid offset: %v", err)
		}
		if sr.Offset < 0 {
			return SearchRequest{}, fmt.Errorf("offset cannot be negative")
		}
	} else {
		sr.Offset = 0
	}

	// Limit
	if s := query.Get("limit"); s != "" {
		if sr.Limit, err = strconv.Atoi(s); err != nil {
			return SearchRequest{}, fmt.Errorf("Invalid limit: %v", err)
		}
		if sr.Limit < 0 {
			return SearchRequest{}, fmt.Errorf("limit cannot be negative")
		}
	} else {
		sr.Limit = maxLimit
	}
	return sr, nil
}
