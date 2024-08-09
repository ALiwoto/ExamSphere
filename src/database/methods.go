package database

import (
	"OnlineExams/src/core/utils/logging"
	"database/sql"
)

func (c *DatabaseContainer) getVersion() (int, error) {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS online_exam_platform_version (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0
	row := c.db.QueryRow("SELECT version FROM online_exam_platform_version LIMIT 1")
	if row != nil {
		_ = row.Scan(&version)
	}
	return version, nil
}

func (c *DatabaseContainer) setVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec("DELETE FROM online_exam_platform_version")
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO online_exam_platform_version (version) VALUES ($1)", version)
	return err
}

func (d *DatabaseContainer) DoMigrations() error {
	version, err := d.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Migrations); version++ {
		var tx *sql.Tx
		tx, err = d.db.Begin()
		if err != nil {
			return err
		}

		migrateFunc := Migrations[version]
		logging.Infof("Migrating to version %d", version+1)
		err = migrateFunc(tx, d)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = d.setVersion(tx, version+1); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

	}

	return nil
}
