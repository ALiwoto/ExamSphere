package database

import (
	"OnlineExams/src/core/utils/logging"
	"context"

	"github.com/jackc/pgx/v5"
)

func (c *DatabaseContainer) getVersion() (int, error) {
	_, err := c.db.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS online_exam_platform_version (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0
	err = c.db.QueryRow(context.Background(),
		"SELECT version FROM online_exam_platform_version LIMIT 1").Scan(&version)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return -1, err
	}
	return version, nil
}

func (c *DatabaseContainer) setVersion(tx pgx.Tx, version int) error {
	_, err := tx.Exec(context.Background(),
		"DELETE FROM online_exam_platform_version")
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO online_exam_platform_version (version) VALUES ($1)", version)
	return err
}

func (d *DatabaseContainer) DoMigrations() error {
	version, err := d.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Migrations); version++ {
		var tx pgx.Tx
		tx, err = d.db.Begin(context.Background())
		if err != nil {
			return err
		}

		migrateFunc := Migrations[version]
		logging.Infof("Migrating to version %d", version+1)
		err = migrateFunc(tx, d)
		if err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}

		if err = d.setVersion(tx, version+1); err != nil {
			return err
		}

		if err = tx.Commit(context.Background()); err != nil {
			return err
		}

	}

	return nil
}
