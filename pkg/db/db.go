package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// schema defines the SQL schema for the scheduler table.
const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
CREATE INDEX idx_scheduler_date ON scheduler (date);`

var db *sql.DB

// Init initializes the database connection and creates the scheduler table if it does not exist.
// It checks for the existence of the database file specified by the TODO_DBFILE environment variable.
// If the file does not exist, it creates the table using the defined schema.
func Init(dbFile string) error {
	var err error
	var install bool

	// Check if the database file exists
	_, err = os.Stat(dbFile)
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
