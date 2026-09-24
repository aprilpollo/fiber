package query

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	defaultLimit = 20

	// MaxLimit caps _limit so one request cannot ask the database for an
	// unbounded number of rows (?_limit=999999999 would otherwise be passed
	// straight through to SQL).
	MaxLimit = 200
)

// Parse builds a QueryOptions from query params (map[string]string, e.g. Fiber's c.Queries()).
//
// Malformed reserved params are rejected instead of being silently ignored, so a
// typo like ?_limit=abc surfaces as a 400 rather than quietly falling back to the
// default page size.
func Parse(table string, params map[string]string) (QueryOptions, error) {
	if !isValidIdentifier(table) {
		return QueryOptions{}, fmt.Errorf("invalid table name: %q", table)
	}

	// Convert map[string]string → map[string][]string for ParseFilters
	multi := make(map[string][]string, len(params))
	for k, v := range params {
		multi[k] = []string{v}
	}

	opts := QueryOptions{
		Table:   table,
		Filters: ParseFilters(multi),
		Limit:   defaultLimit,
		Order:   "ASC",
	}

	if v, ok := params["_sort"]; ok {
		if !isValidIdentifier(v) {
			return QueryOptions{}, fmt.Errorf("_sort: invalid column name %q", v)
		}
		opts.Sort = v
	}

	if v, ok := params["_order"]; ok {
		switch strings.ToUpper(v) {
		case "ASC":
			opts.Order = "ASC"
		case "DESC":
			opts.Order = "DESC"
		default:
			return QueryOptions{}, fmt.Errorf("_order: must be ASC or DESC, got %q", v)
		}
	}

	if v, ok := params["_limit"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return QueryOptions{}, fmt.Errorf("_limit: must be an integer, got %q", v)
		}
		if n < 1 || n > MaxLimit {
			return QueryOptions{}, fmt.Errorf("_limit: must be between 1 and %d, got %d", MaxLimit, n)
		}
		opts.Limit = n
	}

	if v, ok := params["_offset"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return QueryOptions{}, fmt.Errorf("_offset: must be an integer, got %q", v)
		}
		if n < 0 {
			return QueryOptions{}, fmt.Errorf("_offset: must not be negative, got %d", n)
		}
		opts.Offset = n
	}

	// _page overrides _offset. Resolved after _limit so it uses the final page size.
	if v, ok := params["_page"]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return QueryOptions{}, fmt.Errorf("_page: must be an integer, got %q", v)
		}
		if n < 1 {
			return QueryOptions{}, fmt.Errorf("_page: must be 1 or greater, got %d", n)
		}
		if n-1 > math.MaxInt/opts.Limit {
			return QueryOptions{}, fmt.Errorf("_page: %d is too large", n)
		}
		opts.Offset = (n - 1) * opts.Limit
	}

	return opts, nil
}

// Build converts QueryOptions into a parameterized SQL query and args slice.
// Uses $1, $2, ... placeholders (PostgreSQL style).
func Build(opts QueryOptions) (query string, args []interface{}, err error) {
	if opts.Table == "" {
		return "", nil, fmt.Errorf("table is required")
	}

	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	sb.WriteString(opts.Table)

	var conditions []string
	argIdx := 1

	for _, f := range opts.Filters {
		if !isValidIdentifier(f.Field) {
			return "", nil, fmt.Errorf("invalid field name: %q", f.Field)
		}

		switch f.Operator {
		case "IS NULL", "IS NOT NULL":
			conditions = append(conditions, fmt.Sprintf("%s %s", f.Field, f.Operator))

		case "IN", "NOT IN":
			vals, ok := f.Value.([]string)
			if !ok || len(vals) == 0 {
				continue
			}
			placeholders := make([]string, len(vals))
			for i, v := range vals {
				placeholders[i] = fmt.Sprintf("$%d", argIdx)
				args = append(args, v)
				argIdx++
			}
			conditions = append(conditions,
				fmt.Sprintf("%s %s (%s)", f.Field, f.Operator, strings.Join(placeholders, ", ")))

		default:
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", f.Field, f.Operator, argIdx))
			args = append(args, f.Value)
			argIdx++
		}
	}

	if len(conditions) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(conditions, " AND "))
	}

	if opts.Sort != "" {
		sb.WriteString(fmt.Sprintf(" ORDER BY %s %s", opts.Sort, opts.Order))
	}

	sb.WriteString(fmt.Sprintf(" LIMIT $%d", argIdx))
	args = append(args, opts.Limit)
	argIdx++

	sb.WriteString(fmt.Sprintf(" OFFSET $%d", argIdx))
	args = append(args, opts.Offset)

	return sb.String(), args, nil
}
