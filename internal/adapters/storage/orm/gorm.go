package orm

import (
	"aprilpollo/internal/adapters/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
	Close() error
	Migrate() error
}

type GormDB struct {
	db *gorm.DB
}

func NewGormDB(config *config.Database, gormConfig *gorm.Config, isMigration ...bool) (*GormDB, error) {
	var uri string
	if len(isMigration) > 0 && isMigration[0] {
		uri = config.URI_MIGRATION
	} else {
		uri = config.URI
	}

	db, err := gorm.Open(postgres.Open(uri), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	return &GormDB{db}, nil
}

func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}

func (g *GormDB) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (g *GormDB) Migrate(dsls interface{}) error {
	return g.db.AutoMigrate(dsls)
}
