package query

import (
	"fmt"
	"strconv"
	"strings"
)

const defaultLimit = 20

// Parse builds a QueryOptions from Fiber query params (map[string]string via c.Queries())
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

	if v, ok := params["_sort"]; ok && isValidIdentifier(v) {
		opts.Sort = v
	}
	if v := strings.ToUpper(params["_order"]); v == "DESC" {
		opts.Order = "DESC"
	}
	if v, ok := params["_limit"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			opts.Limit = n
		}
	}
	if v, ok := params["_offset"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			opts.Offset = n
		}
	}
	// _page overrides _offset
	if v, ok := params["_page"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			opts.Offset = (n - 1) * opts.Limit
		}
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
