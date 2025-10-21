package hw3

import (
	"sort"
	"strings"
)

func filterUsers(users []User, filter string) []User {
	filtered := make([]User, 0, len(users))
	if filter == "" {
		filtered = append(filtered, users...)
	} else {
		for _, u := range users {
			if strings.Contains(u.Name, filter) || strings.Contains(u.About, filter) {
				filtered = append(filtered, u)
			}
		}
	}
	return filtered
}

func sortUsers(users []User, orderField string, orderBy int) {
	switch {
	case orderBy == 1 && orderField == "id":
		sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	case orderBy == 1 && orderField == "age":
		sort.Slice(users, func(i, j int) bool { return users[i].Age < users[j].Age })
	case orderBy == 1 && orderField == "name":
		sort.Slice(users, func(i, j int) bool { return users[i].Name < users[j].Name })
	case orderBy == -1 && orderField == "id":
		sort.Slice(users, func(i, j int) bool { return users[i].ID > users[j].ID })
	case orderBy == -1 && orderField == "age":
		sort.Slice(users, func(i, j int) bool { return users[i].Age > users[j].Age })
	case orderBy == -1 && orderField == "name":
		sort.Slice(users, func(i, j int) bool { return users[i].Name > users[j].Name })
	}
}

func paginateUsers(users []User, offset int, limit int) []User {
	offset, limit = max(0, offset), max(0, limit) // non-negative additional check
	limit = min(limit, maxLimit)                  // max limit additional check

	offset = min(offset, len(users))
	end := offset + limit
	end = min(end, len(users))

	return users[offset:end]
}
