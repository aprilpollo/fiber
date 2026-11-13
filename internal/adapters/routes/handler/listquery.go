package handler

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"aprilpollo/internal/core/domain"
)

const (
	defaultLimit = 20

	// MaxLimit caps _limit so one request cannot ask the database for an
	// unbounded number of rows (?_limit=999999999 would otherwise be passed
	// straight through to SQL).
	MaxLimit = 200
)

// reserved params that are not filters
var reservedParams = map[string]bool{
	"_limit": true, "_offset": true, "_sort": true, "_order": true, "_page": true,
}

// Order matters: longer suffixes must be checked first to avoid partial match
// ("_nin" before "_in", "_notnull" before "_null").
var operatorSuffixes = []domain.Operator{
	domain.OpContains, domain.OpNotNull, domain.OpIsNull,
	domain.OpGte, domain.OpLte, domain.OpGt, domain.OpLt,
	domain.OpNe, domain.OpEq, domain.OpNotIn, domain.OpIn,
}

// parseListQuery turns query-string params (Fiber's c.Queries()) into a
// domain.ListQuery. It only reads the URL syntax; whether a field may be
// filtered or sorted on is decided by the repository's column allowlist.
//
// Malformed reserved params are rejected instead of being silently ignored, so a
// typo like ?_limit=abc surfaces as a 400 rather than quietly falling back to the
// default page size.
func parseListQuery(params map[string]string) (domain.ListQuery, error) {
	q := domain.ListQuery{Limit: defaultLimit}

	// Sorted so a request with several bad filters always reports the same one.
	keys := make([]string, 0, len(params))
	for k := range params {
		if !reservedParams[k] {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	for _, key := range keys {
		if f, ok := parseFilter(key, params[key]); ok {
			q.Filters = append(q.Filters, f)
		}
	}

	if v, ok := params["_sort"]; ok {
		if v == "" {
			return domain.ListQuery{}, fmt.Errorf("_sort: must name a field")
		}
		q.Sort = v
	}

	if v, ok := params["_order"]; ok {
		switch strings.ToUpper(v) {
		case "ASC":
			q.Desc = false
		case "DESC":
			q.Desc = true
		default:
			return domain.ListQuery{}, fmt.Errorf("_order: must be ASC or DESC, got %q", v)
		}
	}

	if v, ok := params["_limit"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return domain.ListQuery{}, fmt.Errorf("_limit: must be an integer, got %q", v)
		}
		if n < 1 || n > MaxLimit {
			return domain.ListQuery{}, fmt.Errorf("_limit: must be between 1 and %d, got %d", MaxLimit, n)
		}
		q.Limit = n
	}

	if v, ok := params["_offset"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return domain.ListQuery{}, fmt.Errorf("_offset: must be an integer, got %q", v)
		}
		if n < 0 {
			return domain.ListQuery{}, fmt.Errorf("_offset: must not be negative, got %d", n)
		}
		q.Offset = n
	}

	// _page overrides _offset. Resolved after _limit so it uses the final page size.
	if v, ok := params["_page"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return domain.ListQuery{}, fmt.Errorf("_page: must be an integer, got %q", v)
		}
		if n < 1 {
			return domain.ListQuery{}, fmt.Errorf("_page: must be 1 or greater, got %d", n)
		}
		if n-1 > math.MaxInt/q.Limit {
			return domain.ListQuery{}, fmt.Errorf("_page: %d is too large", n)
		}
		q.Offset = (n - 1) * q.Limit
	}

	return q, nil
}

// parseFilter reads "name_contains=john" as {name, contains, "john"}. A key
// with no known suffix is an equality filter on the whole key.
func parseFilter(key, raw string) (domain.Filter, bool) {
	field, op := key, domain.OpEq
	for _, s := range operatorSuffixes {
		if tail := "_" + string(s); strings.HasSuffix(key, tail) && len(key) > len(tail) {
			field, op = key[:len(key)-len(tail)], s
			break
		}
	}

	switch op {
	case domain.OpIsNull, domain.OpNotNull:
		return domain.Filter{Field: field, Op: op}, true
	case domain.OpIn, domain.OpNotIn:
		vals := splitCSV(raw)
		if len(vals) == 0 {
			return domain.Filter{}, false
		}
		return domain.Filter{Field: field, Op: op, Value: vals}, true
	default:
		return domain.Filter{Field: field, Op: op, Value: raw}, true
	}
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
