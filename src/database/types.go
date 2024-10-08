package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// underlyingDbType is the underlying type of the database connection.
// made like this for faster switching between types.
type underlyingDbType = *pgxpool.Pool

// Scannable is an interface that represents a type that is
// capable of scanning values from a database query result,
// such as a *sql.Rows or *sql.Row.
type Scannable interface {
	Scan(dest ...interface{}) error
}

// DatabaseContainer is a struct that holds a database connection
// and the dialect of the database.
type DatabaseContainer struct {
	db      underlyingDbType
	dialect string

	DatabaseErrorHandler func(action string, attemptIndex int, err error) (retry bool)
}
