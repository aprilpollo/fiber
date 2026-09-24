package query

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	opts, err := Parse("users", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Table != "users" || opts.Limit != defaultLimit || opts.Offset != 0 || opts.Order != "ASC" || opts.Sort != "" {
		t.Fatalf("unexpected defaults: %+v", opts)
	}
}

func TestParseValid(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]string
		limit, offset int
		sort, order   string
	}{
		{"limit", map[string]string{"_limit": "50"}, 50, 0, "", "ASC"},
		{"max limit accepted", map[string]string{"_limit": strconv.Itoa(MaxLimit)}, MaxLimit, 0, "", "ASC"},
		{"offset", map[string]string{"_offset": "40"}, defaultLimit, 40, "", "ASC"},
		{"offset zero", map[string]string{"_offset": "0"}, defaultLimit, 0, "", "ASC"},
		{"page 1", map[string]string{"_page": "1"}, defaultLimit, 0, "", "ASC"},
		{"page with limit", map[string]string{"_page": "3", "_limit": "10"}, 10, 20, "", "ASC"},
		{"page overrides offset", map[string]string{"_page": "2", "_limit": "10", "_offset": "999"}, 10, 10, "", "ASC"},
		{"sort", map[string]string{"_sort": "created_at"}, defaultLimit, 0, "created_at", "ASC"},
		{"order desc", map[string]string{"_order": "desc"}, defaultLimit, 0, "", "DESC"},
		{"order asc explicit", map[string]string{"_order": "ASC"}, defaultLimit, 0, "", "ASC"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := Parse("users", tc.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.Limit != tc.limit || opts.Offset != tc.offset || opts.Sort != tc.sort || opts.Order != tc.order {
				t.Fatalf("got limit=%d offset=%d sort=%q order=%q, want limit=%d offset=%d sort=%q order=%q",
					opts.Limit, opts.Offset, opts.Sort, opts.Order, tc.limit, tc.offset, tc.sort, tc.order)
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
		{"sort with sql injection", map[string]string{"_sort": "name; DROP TABLE users"}, "_sort"},
		{"sort empty", map[string]string{"_sort": ""}, "_sort"},
		{"order invalid", map[string]string{"_order": "sideways"}, "_order"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("users", tc.params)
			if err == nil {
				t.Fatalf("want error mentioning %s, got nil", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error mentioning %s, got %q", tc.want, err)
			}
		})
	}
}

func TestParseInvalidTable(t *testing.T) {
	if _, err := Parse("users; DROP TABLE users", nil); err == nil {
		t.Fatal("want error for invalid table name")
	}
}

// Reserved params must never be treated as column filters.
func TestParseReservedParamsAreNotFilters(t *testing.T) {
	opts, err := Parse("users", map[string]string{
		"_limit": "10", "_offset": "0", "_page": "1", "_sort": "name", "_order": "DESC",
		"email_contains": "john",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.Filters) != 1 {
		t.Fatalf("want exactly 1 filter, got %d: %+v", len(opts.Filters), opts.Filters)
	}
	f := opts.Filters[0]
	if f.Field != "email" || f.Operator != "LIKE" || f.Value != "%john%" {
		t.Fatalf("unexpected filter: %+v", f)
	}
}
