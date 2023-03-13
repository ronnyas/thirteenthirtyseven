package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupDatabaseScheme(db *sql.DB) error {
	_, err := db.Exec(`
		create table if not exists 
		points (
			id integer not null primary key,
			timestamp text,
			user_id text, 
			points integer
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
