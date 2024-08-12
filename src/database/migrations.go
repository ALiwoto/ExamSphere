package database

import (
	"ExamSphere/src/core/utils/logging"
	"ExamSphere/src/database/dbScripts"
	"context"

	"github.com/jackc/pgx/v5"
)

func migrateV1(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(context.Background(),
		dbScripts.Migration1Str)
	if err != nil {
		logging.Error("migrateV1: Failed to execute migration 1: ", err)
		return err
	}

	return nil
}

func migrateV2(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(context.Background(),
		dbScripts.Migration2Str)
	if err != nil {
		return err
	}

	return nil
}

func migrateV3(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(context.Background(),
		dbScripts.Migration3Str)
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
