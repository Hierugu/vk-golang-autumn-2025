package hw3

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindUsersIntegration(t *testing.T) {
	db = XMLDatabase{filePath: "dataset.xml"}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	client := SearchClient{AccessToken: accessKey, URL: ts.URL}

	t.Run("BasicLimitAndNextPage", func(t *testing.T) {
		req := SearchRequest{Limit: 1, Offset: 0, Query: "", OrderField: "", OrderBy: 0}
		res, err := client.FindUsers(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Users) > req.Limit {
			t.Fatalf("got more users than limit: %d > %d", len(res.Users), req.Limit)
		}
		if !res.NextPage {
			t.Fatalf("expected NextPage true for limit=1, got false")
		}
	})

	t.Run("FilterByName", func(t *testing.T) {
		req := SearchRequest{Limit: 10, Offset: 0, Query: "Boyd", OrderField: "", OrderBy: 0}
		res, err := client.FindUsers(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Users) == 0 {
			t.Fatalf("expected at least one user matching 'Boyd', got 0")
		}
		found := false
		for _, u := range res.Users {
			if u.ID == 0 {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected user with ID 0 (Boyd) in results")
		}
	})

	t.Run("BadAccessToken", func(t *testing.T) {
		bad := SearchClient{AccessToken: "wrong-token", URL: ts.URL}
		_, err := bad.FindUsers(SearchRequest{Limit: 1})
		if err == nil {
			t.Fatalf("expected error with bad token, got nil")
		}
	})

	t.Run("NotGETRequest", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("AccessToken", accessKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})

	t.Run("InvalidOrderField", func(t *testing.T) {
		_, err := client.FindUsers(SearchRequest{Limit: 1, OrderField: "invalid_field"})
		if err == nil {
			t.Fatalf("expected error for invalid order_field, got nil")
		}
	})

	t.Run("DBErrorReturns500", func(t *testing.T) {
		db = XMLDatabase{filePath: "nonexistent-dataset.xml"}
		_, err := client.FindUsers(SearchRequest{Limit: 1})
		if err == nil {
			t.Fatalf("expected error when DB cannot be opened, got nil")
		}
	})
}

func TestUsersSorting(t *testing.T) {
	users := []User{
		{ID: 2, Name: "Ann", Age: 30},
		{ID: 1, Name: "Bob", Age: 25},
		{ID: 4, Name: "Carl", Age: 20},
		{ID: 3, Name: "Zed", Age: 40},
	}

	checkOrder := func(t *testing.T, got []User, wantIDs []int) {
		if len(got) != len(wantIDs) {
			t.Fatalf("got length %d, want %d", len(got), len(wantIDs))
		}
		for i, u := range got {
			if u.ID != wantIDs[i] {
				t.Fatalf("position %d: got ID %d, want ID %d", i, u.ID, wantIDs[i])
			}
		}
	}

	t.Run("SortByNameAsc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "name", 1)
		want := []int{2, 1, 4, 3} // Ann, Bob, Carl, Zed
		checkOrder(t, in, want)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "name", -1)
		want := []int{3, 4, 1, 2} // Zed, Carl, Bob, Ann
		checkOrder(t, in, want)
	})

	t.Run("SortByAgeAsc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "age", 1)
		want := []int{4, 1, 2, 3} // ages 20, 25, 30, 40
		checkOrder(t, in, want)
	})

	t.Run("SortByAgeDesc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "age", -1)
		want := []int{3, 2, 1, 4} // ages 40, 30, 25, 20
		checkOrder(t, in, want)
	})

	t.Run("SortByIdAsc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "id", 1)
		want := []int{1, 2, 3, 4} // IDs in ascending order
		checkOrder(t, in, want)
	})

	t.Run("SortByIdDesc", func(t *testing.T) {
		in := append([]User(nil), users...)
		sortUsers(in, "id", -1)
		want := []int{4, 3, 2, 1} // IDs in descending order
		checkOrder(t, in, want)
	})

	t.Run("SortByIdWithDuplicates", func(t *testing.T) {
		dupes := []User{
			{ID: 2, Name: "Ann", Age: 30},
			{ID: 1, Name: "Bob", Age: 25},
			{ID: 2, Name: "Carl", Age: 20},
			{ID: 1, Name: "Zed", Age: 40},
		}
		sortUsers(dupes, "id", 1)
		want := []int{1, 1, 2, 2}
		checkOrder(t, dupes, want)
	})

	t.Run("SortByIdEmptySlice", func(t *testing.T) {
		var empty []User
		sortUsers(empty, "id", 1)
		if len(empty) != 0 {
			t.Fatalf("expected empty slice, got %d elements", len(empty))
		}
	})
}

func TestRequestParsing(t *testing.T) {
	t.Run("OrderByNotNumber", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?order_by=abc", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for non-numeric order_by, got nil")
		}
	})

	t.Run("OrderByValidValues", func(t *testing.T) {
		for _, val := range []string{"-1", "0", "1"} {
			req, _ := http.NewRequest("GET", "/?order_by="+val, nil)
			parsed, err := parseQuery(req)
			if err != nil {
				t.Fatalf("unexpected error for order_by=%s: %v", val, err)
			}
			want := 0
			if val == "-1" {
				want = -1
			} else if val == "1" {
				want = 1
			}
			if parsed.OrderBy != want {
				t.Fatalf("order_by: got %d, want %d", parsed.OrderBy, want)
			}
		}
	})

	t.Run("OrderByOutOfRange", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?order_by=234", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for out-of-range order_by, got nil")
		}
	})

	t.Run("OffsetNotNumber", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?offset=abc", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for non-numeric offset, got nil")
		}
	})

	t.Run("OffsetNegative", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?offset=-5", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for negative offset, got nil")
		}
	})

	t.Run("OffsetMissing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		parsed, err := parseQuery(req)
		if err != nil {
			t.Fatalf("unexpected error for missing offset: %v", err)
		}
		if parsed.Offset != 0 {
			t.Fatalf("expected default offset 0, got %d", parsed.Offset)
		}
	})

	t.Run("LimitNotNumber", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?limit=abc", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for non-numeric limit, got nil")
		}
	})

	t.Run("LimitNegative", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?limit=-10", nil)
		_, err := parseQuery(req)
		if err == nil {
			t.Fatalf("expected error for negative limit, got nil")
		}
	})

	t.Run("LimitMissing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		parsed, err := parseQuery(req)
		if err != nil {
			t.Fatalf("unexpected error for missing limit: %v", err)
		}
		if parsed.Limit != 25 {
			t.Fatalf("expected default limit 25, got %d", parsed.Limit)
		}
	})
}

func TestFindUsers(t *testing.T) {
	t.Run("LimitNegative", func(t *testing.T) {
		req := SearchRequest{Limit: -1}
		_, err := (&SearchClient{AccessToken: accessKey, URL: "http://localhost"}).FindUsers(req)
		if err == nil {
			t.Fatalf("expected error for negative limit, got nil")
		}
	})

	t.Run("LimitTooLarge", func(t *testing.T) {
		req := SearchRequest{Limit: 100}
		_, err := (&SearchClient{AccessToken: accessKey, URL: "http://localhost"}).FindUsers(req)
		if err == nil {
			t.Fatalf("expected error for limit > 25, got nil")
		}
	})

	t.Run("OffsetNegative", func(t *testing.T) {
		req := SearchRequest{Offset: -5}
		_, err := (&SearchClient{AccessToken: accessKey, URL: "http://localhost"}).FindUsers(req)
		if err == nil {
			t.Fatalf("expected error for negative offset, got nil")
		}
	})

	t.Run("ServerNotReachable", func(t *testing.T) {
		client := SearchClient{AccessToken: accessKey, URL: "http://127.0.0.1:9999"}
		_, err := client.FindUsers(SearchRequest{Limit: 1})
		if err == nil {
			t.Fatalf("expected error when server is not reachable, got nil")
		}
	})

	t.Run("BadRequestCantUnpackErrorJSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{invalid json"))
		}))
		defer ts.Close()
		client := SearchClient{AccessToken: accessKey, URL: ts.URL}
		_, err := client.FindUsers(SearchRequest{Limit: 1})
		if err == nil || err.Error() == "" {
			t.Fatalf("expected error for bad request with invalid error json, got nil")
		}
	})

	t.Run("BadRequestUnknownError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"unknown bad request error"}`))
		}))
		defer ts.Close()
		client := SearchClient{AccessToken: accessKey, URL: ts.URL}
		_, err := client.FindUsers(SearchRequest{Limit: 1})
		if err == nil || err.Error() == "" {
			t.Fatalf("expected error for unknown bad request error, got nil")
		}
	})

	t.Run("StatusOKCantUnpackResultJSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json"))
		}))
		defer ts.Close()
		client := SearchClient{AccessToken: accessKey, URL: ts.URL}
		_, err := client.FindUsers(SearchRequest{Limit: 1})
		if err == nil {
			t.Fatalf("expected error for invalid result json, got nil")
		}
	})
}
