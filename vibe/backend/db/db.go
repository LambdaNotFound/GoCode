package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Init(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			amount      REAL    NOT NULL,
			category    TEXT    NOT NULL,
			description TEXT    NOT NULL DEFAULT '',
			date        TEXT    NOT NULL,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
