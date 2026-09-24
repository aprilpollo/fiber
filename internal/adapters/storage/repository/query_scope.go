package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ApplyToGorm chains filters, sort, limit, offset onto a *gorm.DB scope.
// The caller is responsible for calling .Find(), .Scan(), etc. afterward.
//
//	db := querybuilder.ApplyToGorm(db.Model(&User{}), opts)
//	db.Find(&users)
func ApplyToGorm(db *gorm.DB, opts QueryOptions) *gorm.DB {
	for _, f := range opts.Filters {
		db = applyFilter(db, f)
	}

	if opts.Sort != "" {
		order := opts.Sort
		if strings.ToUpper(opts.Order) == "DESC" {
			order += " DESC"
		}
		db = db.Order(order)
	}

	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		db = db.Offset(opts.Offset)
	}

	return db
}

func applyFilter(db *gorm.DB, f Filter) *gorm.DB {
	col := f.Field // already validated as safe identifier in ParseFilters

	switch f.Operator {
	case "IS NULL":
		return db.Where(fmt.Sprintf("%s IS NULL", col))

	case "IS NOT NULL":
		return db.Where(fmt.Sprintf("%s IS NOT NULL", col))

	case "IN":
		// GORM accepts slice natively with "IN ?"
		return db.Where(fmt.Sprintf("%s IN ?", col), f.Value)

	case "NOT IN":
		return db.Where(fmt.Sprintf("%s NOT IN ?", col), f.Value)

	default:
		// "=", "!=", ">", ">=", "<", "<=", "LIKE"
		return db.Where(fmt.Sprintf("%s %s ?", col, f.Operator), f.Value)
	}
}
