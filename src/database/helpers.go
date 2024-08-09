package database

import (
	"OnlineExams/src/core/appConfig"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
)

// StartDatabase connects to the database and sets the DefaultContainer
// variable to the connected database.
// This function will also do all the required migrations.
func StartDatabase() error {
	container, err := New(dbType, appConfig.GetDBUrl())
	if err != nil {
		return fmt.Errorf("failed to start database: %w", err)
	}

	DefaultContainer = container
	return nil
}

// New connects to the given PostgreSQL database and returns a DatabaseContainer.
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
