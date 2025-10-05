package models

import (
	"time"

	"gorm.io/gorm"
)

type TableNamer interface {
	TableName() string
}

type ModelList []TableNamer

func All() ModelList {
	return ModelList{
		&User{},
		&UserRole{},
	}
}

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
