package models

import (
	"gorm.io/gorm"
	"time"
)

type ModelList []interface{}

func All() ModelList {
	return ModelList{
		&User{},
	}
}

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
