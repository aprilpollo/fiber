package repository

import (
	"errors"
	"reflect"
	"testing"

	"aprilpollo/internal/core/domain"
)

func TestFilterClause(t *testing.T) {
	tests := []struct {
		f        domain.Filter
		wantSQL  string
		wantArgs []any
	}{
		{domain.Filter{Field: "n", Op: domain.OpEq, Value: "x"}, "col = ?", []any{"x"}},
		{domain.Filter{Field: "n", Op: domain.OpNe, Value: "x"}, "col != ?", []any{"x"}},
		{domain.Filter{Field: "n", Op: domain.OpGte, Value: "3"}, "col >= ?", []any{"3"}},
		{domain.Filter{Field: "n", Op: domain.OpContains, Value: "jo"}, "col LIKE ?", []any{"%jo%"}},
		{domain.Filter{Field: "n", Op: domain.OpContains, Value: `50%_\`}, "col LIKE ?", []any{`%50\%\_\\%`}},
		{domain.Filter{Field: "n", Op: domain.OpIn, Value: []string{"a", "b"}}, "col IN ?", []any{[]string{"a", "b"}}},
		{domain.Filter{Field: "n", Op: domain.OpNotIn, Value: []string{"a"}}, "col NOT IN ?", []any{[]string{"a"}}},
		{domain.Filter{Field: "n", Op: domain.OpIsNull}, "col IS NULL", nil},
		{domain.Filter{Field: "n", Op: domain.OpNotNull}, "col IS NOT NULL", nil},
	}
	for _, tc := range tests {
		sql, args, err := filterClause("col", tc.f)
		if err != nil {
			t.Fatalf("%v: %v", tc.f, err)
		}
		if sql != tc.wantSQL || !reflect.DeepEqual(args, tc.wantArgs) {
			t.Errorf("%v: got %q %v, want %q %v", tc.f, sql, args, tc.wantSQL, tc.wantArgs)
		}
	}
}

func TestListScopesRefusesFieldsOutsideTheAllowlist(t *testing.T) {
	cols := columns{"name": "name"}
	for _, q := range []domain.ListQuery{
		{Filters: []domain.Filter{{Field: "password_hash", Op: domain.OpEq, Value: "x"}}},
		{Filters: []domain.Filter{{Field: "name; DROP TABLE users", Op: domain.OpEq, Value: "x"}}},
		{Sort: "password_hash"},
		{Filters: []domain.Filter{{Field: "name", Op: "regex", Value: "x"}}},
	} {
		_, _, err := listScopes(q, cols)
		var verr *domain.ValidationError
		if !errors.As(err, &verr) {
			t.Errorf("%+v: want ValidationError, got %v", q, err)
		}
	}
}
