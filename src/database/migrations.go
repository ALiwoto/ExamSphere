package database

import (
	"ExamSphere/src/core/utils/logging"
	"ExamSphere/src/database/dbScripts"

	"github.com/jackc/pgx/v5"
)

func migrateV1(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(container.MigrationCtx, dbScripts.Migration1Str)
	if err != nil {
		logging.Error("migrateV1: Failed to execute migration 1: ", err)
		return err
	}

	return nil
}

func migrateV2(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(container.MigrationCtx, dbScripts.Migration2Str)
	if err != nil {
		logging.Error("migrateV2: Failed to execute migration 2: ", err)
		return err
	}

	return nil
}

func migrateV3(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(container.MigrationCtx, dbScripts.Migration3Str)
	if err != nil {
		logging.Error("migrateV3: Failed to execute migration 3: ", err)
		return err
	}

	return nil
}

func migrateV4(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(container.MigrationCtx, dbScripts.Migration4Str)
	if err != nil {
		logging.Error("migrateV4: Failed to execute migration 4: ", err)
		return err
	}

	return nil
}

func migrateV5(tx pgx.Tx, container *DatabaseContainer) error {
	_, err := tx.Exec(container.MigrationCtx, dbScripts.Migration5Str)
	if err != nil {
		logging.Error("migrateV5: Failed to execute migration 5: ", err)
		return err
	}

	return nil
}
