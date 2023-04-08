package game

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type leaderboardConfig struct {
	name    string
	sqlStmt string
	prefix  string
}

var Config struct {
	mainChannel string
	db 			*sql.DB
	StreakDays  int
}

func SetDatabase(db *sql.DB) {
	Config.db = db
}
func SetMainChannel(channelID string) {
	Config.mainChannel = channelID
}
func SetStreakDays(days int) {
	Config.StreakDays = days
}

func calculatePointsFromTimestamp(timestamp time.Time) int {
	layout := "2006-01-02 15:04:05.000 -0700 MST"

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	timestamp, err = time.ParseInLocation(layout, timestamp.Format(layout), loc)
	if err != nil {
		panic(err)
	}

	points := 60 - timestamp.Second()
	
	return points
}

func generateLeaderboardMessage(prefix string, rows *sql.Rows) (string, error) {
	leaderboardMessage := prefix
	rank := 1

	for rows.Next() {
		var userID string
		var points int

		if err := rows.Scan(&userID, &points); err != nil {
			return "", err
		}

		switch rank {
			case 1:
				leaderboardMessage += ":first_place: "
			case 2:
				leaderboardMessage += ":second_place: "
			case 3:
				leaderboardMessage += ":third_place: "
			default:
				leaderboardMessage += ":medal: "
		}

		leaderboardMessage += fmt.Sprintf("%s %d\n", userID, points)
		rank++
	}

	return leaderboardMessage, nil
}

func SavePoints(userID string, points int) bool {
	db := Config.db
	sqlStmt := `
		select count(*) from points
		where user_id = ? and timestamp >= date('now');
		`
	var count int
	err := db.QueryRow(sqlStmt, userID).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false
	}
	
	if count > 0 {
		return false
	} 
		
	sqlStmt = `
	insert into points (timestamp, user_id, points) 
	values (?, ?, ?);
	`
	_, err = db.Exec(sqlStmt, time.Now(), userID, points)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}