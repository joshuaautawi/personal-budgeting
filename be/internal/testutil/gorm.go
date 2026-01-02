package testutil

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"personal-budgeting/be/internal/dbmodel"
)

func NewTestGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := dbmodel.AutoMigrate(db); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	// Enable FK constraints in SQLite so behavior is closer to Postgres.
	_ = db.Exec("PRAGMA foreign_keys = ON").Error
	return db
}


