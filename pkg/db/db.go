package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128),
	comment TEXT,
	repeat VARCHAR(128)
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	fmt.Println("Initializing DB with file:", dbFile)

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	} else {
		row := DB.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler'`)
		var name string
		err = row.Scan(&name)
		if err == sql.ErrNoRows || name != "scheduler" {
			_, err = DB.Exec(schema)
			if err != nil {
				return fmt.Errorf("failed to create schema on existing DB: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check existing schema: %w", err)
		}
	}

	return nil
}

func Close() error {
	return DB.Close()
}
