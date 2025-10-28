package main

import (
	"encoding/json"
	"net/http"
)

const accessKey = "12547890"
const maxLimit = 25

var db Database = XMLDatabase{"dataset.xml"}
var allowedOrderFields = []string{"id", "age", "name", ""}
var allowedOrderByFields = []int{-1, 0, 1}

// // utils
// func filterUsers(users []User, filter string) []User {
// 	filtered := make([]User, 0, len(users))
// 	if filter == "" {
// 		filtered = append(filtered, users...)
// 	} else {
// 		for _, u := range users {
// 			if strings.Contains(u.Name, filter) || strings.Contains(u.About, filter) {
// 				filtered = append(filtered, u)
// 			}
// 		}
// 	}
// 	return filtered
// }

// func sortUsers(users []User, orderField string, orderBy int) {
// 	switch {
// 	case orderBy == 1 && orderField == "id":
// 		sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
// 	case orderBy == 1 && orderField == "age":
// 		sort.Slice(users, func(i, j int) bool { return users[i].Age < users[j].Age })
// 	case orderBy == 1 && orderField == "name":
// 		sort.Slice(users, func(i, j int) bool { return users[i].Name < users[j].Name })
// 	case orderBy == -1 && orderField == "id":
// 		sort.Slice(users, func(i, j int) bool { return users[i].ID > users[j].ID })
// 	case orderBy == -1 && orderField == "age":
// 		sort.Slice(users, func(i, j int) bool { return users[i].Age > users[j].Age })
// 	case orderBy == -1 && orderField == "name":
// 		sort.Slice(users, func(i, j int) bool { return users[i].Name > users[j].Name })
// 	}
// }

// func paginateUsers(users []User, offset int, limit int) []User {
// 	offset, limit = max(0, offset), max(0, limit) // non-negative additional check
// 	limit = min(limit, maxLimit)                  // max limit additional check

// 	offset = min(offset, len(users))
// 	end := offset + limit
// 	end = min(end, len(users))

// 	return users[offset:end]
// }

// // db
// type XMLUser struct {
// 	ID        int    `xml:"id"`
// 	FirstName string `xml:"first_name"`
// 	LastName  string `xml:"last_name"`
// 	Age       int    `xml:"age"`
// 	About     string `xml:"about"`
// 	Gender    string `xml:"gender"`
// }

// func (u XMLUser) ToUser() User {
// 	return User{u.ID, u.FirstName + " " + u.LastName, u.Age, u.About, u.Gender}
// }

// type Database interface {
// 	Load() ([]User, error)
// }

// type XMLDatabase struct {
// 	filePath string
// }

// func (db XMLDatabase) Load() ([]User, error) {
// 	xmlData, err := os.Open(db.filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer xmlData.Close()

// 	rawUser := XMLUser{}
// 	users := []User{}
// 	d := xml.NewDecoder(xmlData)
// 	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
// 		switch se := t.(type) {
// 		case xml.StartElement:
// 			if se.Name.Local == "row" {
// 				d.DecodeElement(&rawUser, &se)
// 				users = append(users, rawUser.ToUser())
// 			}
// 		}
// 	}
// 	return users, nil
// }

// // parseRequest
// func checkHeader(r *http.Request) bool {
// 	return r.Header.Get("AccessToken") == accessKey
// }

// func checkMethod(r *http.Request) bool {
// 	return r.Method == http.MethodGet
// }

// func parseQuery(r *http.Request) (SearchRequest, error) {
// 	sr := SearchRequest{}
// 	query := r.URL.Query()
// 	var err error

// 	// Query
// 	sr.Query = query.Get("query")

// 	// OrderField
// 	sr.OrderField = query.Get("order_field")
// 	if !slices.Contains(allowedOrderFields, sr.OrderField) {
// 		return SearchRequest{}, fmt.Errorf(ErrorBadOrderField)
// 	}

// 	// OrderBy
// 	if s := query.Get("order_by"); s != "" {
// 		if sr.OrderBy, err = strconv.Atoi(s); err != nil {
// 			return SearchRequest{}, fmt.Errorf("Invalid order_by: %v", err)
// 		}
// 		if !slices.Contains(allowedOrderByFields, sr.OrderBy) {
// 			return SearchRequest{}, fmt.Errorf("Invalid order_by: must be one of -1, 0, 1")
// 		}
// 	} else {
// 		sr.OrderBy = OrderByAsIs
// 	}

// 	// Offset
// 	if s := query.Get("offset"); s != "" {
// 		if sr.Offset, err = strconv.Atoi(s); err != nil {
// 			return SearchRequest{}, fmt.Errorf("Invalid offset: %v", err)
// 		}
// 		if sr.Offset < 0 {
// 			return SearchRequest{}, fmt.Errorf("offset cannot be negative")
// 		}
// 	} else {
// 		sr.Offset = 0
// 	}

// 	// Limit
// 	if s := query.Get("limit"); s != "" {
// 		if sr.Limit, err = strconv.Atoi(s); err != nil {
// 			return SearchRequest{}, fmt.Errorf("Invalid limit: %v", err)
// 		}
// 		if sr.Limit < 0 {
// 			return SearchRequest{}, fmt.Errorf("limit cannot be negative")
// 		}
// 	} else {
// 		sr.Limit = maxLimit
// 	}
// 	return sr, nil
// }

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
