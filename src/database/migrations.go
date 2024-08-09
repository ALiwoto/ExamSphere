package database

import (
	"OnlineExams/src/database/dbScripts"
	"context"

	"github.com/jackc/pgx/v5"
)

func migrateV1(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(context.Background(),
		dbScripts.Migration1Str)
	if err != nil {
		return nil
	}

	return nil
}

// func migrateV2(pgx.Tx, *DatabaseContainer) error {
// 	return nil
// }

// func migrateV3(pgx.Tx, *DatabaseContainer) error {
// 	return nil
// }
