package domain

import (
	"errors"
	"fmt"
)

// Errors the core reports. Adapters translate their own failures into these
// (a missing row becomes ErrNotFound, a unique violation ErrConflict), and the
// HTTP adapter maps them to status codes — so neither side has to know how the
// other one fails.
var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")
)

// ValidationError is input the core refused. Its message is written for the
// caller, so an adapter may show it as is.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func invalid(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}
