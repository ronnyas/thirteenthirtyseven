package database

import (
	"database/sql"
	"time"

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
		special_points (
			id integer not null primary key,
			timestamp text not null,
			user_id text not null,
			points integer not null,
			type text not null
		);

		create table if not exists
		streaks (
			id integer not null primary key,
			user_id text not null,
			start_time text not null,
			end_time text not null
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

// AddPoints adds regular points to the database
func AddPoints(db *sql.DB, userID string, timestamp time.Time, points int) error {
	_, err := db.Exec(`
		INSERT INTO points (timestamp, user_id, points)
		VALUES (?, ?, ?)
	`, timestamp.Format(time.RFC3339), userID, points)

	return err
}

// AddSpecialPoints adds special points to the database with the correct type
func AddSpecialPoints(db *sql.DB, userID string, timestamp time.Time, points int, eggType string) error {
	_, err := db.Exec(`
		INSERT INTO special_points (timestamp, user_id, points, type)
		VALUES (?, ?, ?, ?)
	`, timestamp.Format(time.RFC3339), userID, points, eggType)

	return err
}

// GetTotalPoints combines regular and special points for a user
func GetTotalPoints(db *sql.DB, userID string) (int, error) {
	var regularPoints, specialPoints int

	// Get regular points
	err := db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM points WHERE user_id = ?", userID).Scan(&regularPoints)
	if err != nil {
		return 0, err
	}

	// Get special points
	err = db.QueryRow("SELECT COALESCE(SUM(points), 0) FROM special_points WHERE user_id = ?", userID).Scan(&specialPoints)
	if err != nil {
		return 0, err
	}

	return regularPoints + specialPoints, nil
}

// GetPointsForDay gets the total points for a user on a specific day
func GetPointsForDay(db *sql.DB, userID string, date time.Time) (int, error) {
	var regularPoints, specialPoints int
	dateStr := date.Format("2006-01-02")

	// Get regular points
	err := db.QueryRow(`
		SELECT COALESCE(SUM(points), 0) 
		FROM points 
		WHERE user_id = ? AND date(timestamp) = date(?)
	`, userID, dateStr).Scan(&regularPoints)
	if err != nil {
		return 0, err
	}

	// Get special points
	err = db.QueryRow(`
		SELECT COALESCE(SUM(points), 0) 
		FROM special_points 
		WHERE user_id = ? AND date(timestamp) = date(?)
	`, userID, dateStr).Scan(&specialPoints)
	if err != nil {
		return 0, err
	}

	return regularPoints + specialPoints, nil
}
