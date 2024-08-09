package database

import "database/sql"

var (
	DefaultContainer *DatabaseContainer
)

type MigrationFunc func(*sql.Tx, *DatabaseContainer) error

// Migrations is a list of functions that will migrate a database to the latest version.
var Migrations = [...]MigrationFunc{migrateV1, migrateV2, migrateV3}
