package handler

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"aprilpollo/internal/core/domain"
)

func TestParseDefaults(t *testing.T) {
	q, err := parseListQuery(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if q.Limit != defaultLimit || q.Offset != 0 || q.Desc || q.Sort != "" || len(q.Filters) != 0 {
		t.Fatalf("unexpected defaults: %+v", q)
	}
}

func TestParseValid(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]string
		limit, offset int
		sort          string
		desc          bool
	}{
		{"limit", map[string]string{"_limit": "50"}, 50, 0, "", false},
		{"max limit accepted", map[string]string{"_limit": strconv.Itoa(MaxLimit)}, MaxLimit, 0, "", false},
		{"offset", map[string]string{"_offset": "40"}, defaultLimit, 40, "", false},
		{"offset zero", map[string]string{"_offset": "0"}, defaultLimit, 0, "", false},
		{"page 1", map[string]string{"_page": "1"}, defaultLimit, 0, "", false},
		{"page with limit", map[string]string{"_page": "3", "_limit": "10"}, 10, 20, "", false},
		{"page overrides offset", map[string]string{"_page": "2", "_limit": "10", "_offset": "999"}, 10, 10, "", false},
		{"sort", map[string]string{"_sort": "created_at"}, defaultLimit, 0, "created_at", false},
		{"order desc", map[string]string{"_order": "desc"}, defaultLimit, 0, "", true},
		{"order asc explicit", map[string]string{"_order": "ASC"}, defaultLimit, 0, "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q, err := parseListQuery(tc.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if q.Limit != tc.limit || q.Offset != tc.offset || q.Sort != tc.sort || q.Desc != tc.desc {
				t.Fatalf("got limit=%d offset=%d sort=%q desc=%v, want limit=%d offset=%d sort=%q desc=%v",
					q.Limit, q.Offset, q.Sort, q.Desc, tc.limit, tc.offset, tc.sort, tc.desc)
			}
		})
	}
}

// Malformed reserved params must be rejected, not silently replaced by defaults.
func TestParseRejectsBadInput(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
		want   string
	}{
		{"limit not a number", map[string]string{"_limit": "abc"}, "_limit"},
		{"limit zero", map[string]string{"_limit": "0"}, "_limit"},
		{"limit negative", map[string]string{"_limit": "-5"}, "_limit"},
		{"limit over max", map[string]string{"_limit": strconv.Itoa(MaxLimit + 1)}, "_limit"},
		{"limit absurd", map[string]string{"_limit": "999999999"}, "_limit"},
		{"limit overflows int", map[string]string{"_limit": "99999999999999999999999"}, "_limit"},
		{"offset not a number", map[string]string{"_offset": "abc"}, "_offset"},
		{"offset negative", map[string]string{"_offset": "-1"}, "_offset"},
		{"page not a number", map[string]string{"_page": "abc"}, "_page"},
		{"page zero", map[string]string{"_page": "0"}, "_page"},
		{"page overflows offset", map[string]string{"_page": strconv.Itoa(1 << 60), "_limit": "200"}, "_page"},
		{"sort empty", map[string]string{"_sort": ""}, "_sort"},
		{"order invalid", map[string]string{"_order": "sideways"}, "_order"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseListQuery(tc.params)
			if err == nil {
				t.Fatalf("want error mentioning %s, got nil", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error mentioning %s, got %q", tc.want, err)
			}
		})
	}
}

// Reserved params must never be treated as field filters, and the parser must
// hand the core an operator — not SQL.
func TestParseFilters(t *testing.T) {
	q, err := parseListQuery(map[string]string{
		"_limit": "10", "_offset": "0", "_page": "1", "_sort": "name", "_order": "DESC",
		"email_contains":  "john",
		"name":            "Alice",
		"status_nin":      "a, b,,",
		"deleted_at_null": "true",
		"tags_in":         " , ",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []domain.Filter{
		{Field: "deleted_at", Op: domain.OpIsNull},
		{Field: "email", Op: domain.OpContains, Value: "john"},
		{Field: "name", Op: domain.OpEq, Value: "Alice"},
		{Field: "status", Op: domain.OpNotIn, Value: []string{"a", "b"}},
	}
	if !reflect.DeepEqual(q.Filters, want) {
		t.Fatalf("got  %+v\nwant %+v", q.Filters, want)
	}
}
