package database

import (
	"github.com/jackc/pgx/v5"
)

var (
	DefaultContainer *DatabaseContainer
)

type MigrationFunc func(pgx.Tx, *DatabaseContainer) error

// Migrations is a list of functions that will migrate a database to the latest version.
var Migrations = [...]MigrationFunc{
	migrateV1,
	migrateV2,
	migrateV3,
}
