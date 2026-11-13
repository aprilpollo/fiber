package domain

// Operator is a comparison the core can ask a repository for. It names intent,
// not syntax: turning OpContains into LIKE '%…%' is the storage adapter's job.
type Operator string

const (
	OpEq       Operator = "eq"
	OpNe       Operator = "ne"
	OpGt       Operator = "gt"
	OpGte      Operator = "gte"
	OpLt       Operator = "lt"
	OpLte      Operator = "lte"
	OpContains Operator = "contains"
	OpIn       Operator = "in"
	OpNotIn    Operator = "nin"
	OpIsNull   Operator = "null"
	OpNotNull  Operator = "notnull"
)

// Filter is one condition on an entity field (by its domain name, e.g. "email").
// Value is a string for comparisons, []string for OpIn/OpNotIn and unused for
// the null checks.
type Filter struct {
	Field string
	Op    Operator
	Value any
}

// ListQuery is how a caller asks for a page of entities.
type ListQuery struct {
	Filters []Filter
	Sort    string // domain field name; empty leaves the order to storage
	Desc    bool
	Limit   int
	Offset  int
}
