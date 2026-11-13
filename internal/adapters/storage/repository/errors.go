package repository

import (
	"errors"

	"aprilpollo/internal/core/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// pgUniqueViolation is the SQLSTATE Postgres reports for a duplicate key.
const pgUniqueViolation = "23505"

// translate turns a storage failure into the domain error the core expects,
// leaving anything unrecognised as is for the caller to treat as internal.
func translate(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return domain.ErrConflict
	}
	return err
}
