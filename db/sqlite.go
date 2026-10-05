package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase() (*sql.DB, error) {
	db, error := sql.Open("sqlite3", "./data/cars.db")
	if error != nil {
		return nil, error
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS cars (
		id INTEGER PRIMARY KEY,
		make TEXT NOT NULL,
		model TEXT NOT NULL,
		year INTEGER NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
