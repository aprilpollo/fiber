package repository

import (
	"fmt"
	"strings"

	"aprilpollo/internal/core/domain"

	"gorm.io/gorm"
)

// columns maps a domain field name to the column that stores it. It is the
// allowlist for filtering and sorting: a field missing here is refused before
// any SQL is built, so a caller can neither probe the schema nor reach a column
// the API does not expose.
type columns map[string]string

// listScopes turns a domain.ListQuery into GORM scopes. where narrows the rows
// and is shared by the page query and the count, so total reflects the filters;
// page adds order, limit and offset on top.
func listScopes(q domain.ListQuery, cols columns) (where, page func(*gorm.DB) *gorm.DB, err error) {
	type cond struct {
		sql  string
		args []any
	}
	conds := make([]cond, 0, len(q.Filters))
	for _, f := range q.Filters {
		col, ok := cols[f.Field]
		if !ok {
			return nil, nil, &domain.ValidationError{Field: f.Field, Message: "cannot be filtered on"}
		}
		sql, args, err := filterClause(col, f)
		if err != nil {
			return nil, nil, err
		}
		conds = append(conds, cond{sql, args})
	}

	order := ""
	if q.Sort != "" {
		col, ok := cols[q.Sort]
		if !ok {
			return nil, nil, &domain.ValidationError{Field: q.Sort, Message: "cannot be sorted on"}
		}
		order = col
		if q.Desc {
			order += " DESC"
		}
	}

	where = func(db *gorm.DB) *gorm.DB {
		for _, c := range conds {
			db = db.Where(c.sql, c.args...)
		}
		return db
	}
	page = func(db *gorm.DB) *gorm.DB {
		if order != "" {
			db = db.Order(order)
		}
		if q.Limit > 0 {
			db = db.Limit(q.Limit)
		}
		if q.Offset > 0 {
			db = db.Offset(q.Offset)
		}
		return db
	}
	return where, page, nil
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// filterClause renders one filter as a parameterised condition on col.
// col comes from the columns allowlist, never from the caller.
func filterClause(col string, f domain.Filter) (string, []any, error) {
	switch f.Op {
	case domain.OpIsNull:
		return col + " IS NULL", nil, nil
	case domain.OpNotNull:
		return col + " IS NOT NULL", nil, nil
	case domain.OpIn, domain.OpNotIn:
		vals, ok := f.Value.([]string)
		if !ok || len(vals) == 0 {
			return "", nil, &domain.ValidationError{Field: f.Field, Message: "needs at least one value"}
		}
		op := "IN"
		if f.Op == domain.OpNotIn {
			op = "NOT IN"
		}
		return fmt.Sprintf("%s %s ?", col, op), []any{vals}, nil
	case domain.OpContains:
		// Escape the wildcards so ?name_contains=50% matches a literal "50%".
		return col + " LIKE ?", []any{"%" + likeEscaper.Replace(fmt.Sprint(f.Value)) + "%"}, nil
	}

	sqlOp, ok := map[domain.Operator]string{
		domain.OpEq: "=", domain.OpNe: "!=",
		domain.OpGt: ">", domain.OpGte: ">=",
		domain.OpLt: "<", domain.OpLte: "<=",
	}[f.Op]
	if !ok {
		return "", nil, &domain.ValidationError{Field: f.Field, Message: fmt.Sprintf("unsupported operator %q", f.Op)}
	}
	return fmt.Sprintf("%s %s ?", col, sqlOp), []any{f.Value}, nil
}
