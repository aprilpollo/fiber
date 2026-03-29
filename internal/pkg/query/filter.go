package query

import (
	"strings"
)

// Operator suffixes → SQL operators
// key = suffix after last underscore, value = SQL operator
var suffixToSQL = map[string]string{
	"eq":       "=",
	"ne":       "!=",
	"gt":       ">",
	"gte":      ">=",
	"lt":       "<",
	"lte":      "<=",
	"contains": "LIKE",
	"in":       "IN",
	"nin":      "NOT IN",
	"null":     "IS NULL",
	"notnull":  "IS NOT NULL",
}

// Order matters: longer suffixes must be checked first to avoid partial match
var knownSuffixes = []string{
	"contains", "notnull", "null", "gte", "lte", "gt", "lt", "ne", "eq", "nin", "in",
}

// Filter represents one parsed condition
type Filter struct {
	Field    string      // column name
	Operator string      // SQL operator e.g. "=", "LIKE", "IN"
	Value    interface{} // nil for IS NULL / IS NOT NULL
}

// QueryOptions holds all parsed options from query params
type QueryOptions struct {
	Table   string
	Filters []Filter
	Sort    string
	Order   string // "ASC" or "DESC"
	Limit   int
	Offset  int
}

// reserved params that are not filters
var reservedParams = map[string]bool{
	"_limit": true, "_offset": true, "_sort": true, "_order": true, "_page": true,
}

// ParseFilters extracts filter conditions from query params
func ParseFilters(params map[string][]string) []Filter {
	var filters []Filter

	for key, values := range params {
		if reservedParams[key] || len(values) == 0 {
			continue
		}

		raw := values[0]
		field, suffix := splitFieldSuffix(key)

		if !isValidIdentifier(field) {
			continue
		}

		filter := buildFilter(field, suffix, raw)
		if filter != nil {
			filters = append(filters, *filter)
		}
	}

	return filters
}

func buildFilter(field, suffix, raw string) *Filter {
	switch suffix {
	case "null":
		return &Filter{Field: field, Operator: "IS NULL", Value: nil}

	case "notnull":
		return &Filter{Field: field, Operator: "IS NOT NULL", Value: nil}

	case "contains":
		return &Filter{Field: field, Operator: "LIKE", Value: "%" + raw + "%"}

	case "in":
		vals := splitCSV(raw)
		if len(vals) == 0 {
			return nil
		}
		return &Filter{Field: field, Operator: "IN", Value: vals}

	case "nin":
		vals := splitCSV(raw)
		if len(vals) == 0 {
			return nil
		}
		return &Filter{Field: field, Operator: "NOT IN", Value: vals}

	default:
		sqlOp, ok := suffixToSQL[suffix]
		if !ok {
			sqlOp = "="
		}
		return &Filter{Field: field, Operator: sqlOp, Value: raw}
	}
}

// splitFieldSuffix splits "name_contains" → ("name", "contains")
// if no known suffix, treat entire key as field with "eq"
func splitFieldSuffix(key string) (field, suffix string) {
	for _, s := range knownSuffixes {
		tail := "_" + s
		if strings.HasSuffix(key, tail) {
			return key[:len(key)-len(tail)], s
		}
	}
	return key, "eq"
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// isValidIdentifier allows only [a-zA-Z0-9_] to prevent SQL injection
func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
