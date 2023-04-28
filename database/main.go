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

func SetupDatabaseSchema(db *sql.DB) error {
	_, err := db.Exec(`
		create table if not exists 
		points (
			id integer not null primary key,
			timestamp text not null,
			user_id text not null, 
			points integer not null
		);

		create table if not exists
		streaks (
			id integer not null primary key,
			user_id text not null,
			start_time text not null,
			end_time text not null
		);

		CREATE TABLE IF NOT EXISTS
		config (
			id integer NOT NULL primary key,
			serverid text NOT NULL,
			name text NOT NULL,
			value text NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
