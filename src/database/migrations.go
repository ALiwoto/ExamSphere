package database

import (
	"OnlineExams/src/database/dbScripts"
	"database/sql"
)

func migrateV1(tx *sql.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(dbScripts.Migration1Str)
	if err != nil {
		return nil
	}
	return nil
}

func migrateV2(*sql.Tx, *DatabaseContainer) error {
	return nil
}

func migrateV3(*sql.Tx, *DatabaseContainer) error {
	return nil
}
