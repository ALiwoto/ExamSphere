package database

import (
	"database/sql"
	"fmt"
)

// New connects to the given SQL database and returns a DatabaseContainer.
//
// Only SQLite and Postgres are currently fully supported.
//
// When using SQLite, make sure to enable foreign keys by adding `?_foreign_keys=true`:
//
//	container, err := database.New("sqlite3", "file:database_file.db?_foreign_keys=on", nil)
func New(dialect, address string) (*DatabaseContainer, error) {
	db, err := sql.Open(dialect, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	container := NewWithDB(db, dialect)
	err = container.DoMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return container, nil
}

func NewWithDB(db *sql.DB, dialect string) *DatabaseContainer {
	return &DatabaseContainer{
		db:      db,
		dialect: dialect,
	}
}
