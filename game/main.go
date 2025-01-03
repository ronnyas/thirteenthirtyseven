package game

import (
	"database/sql"
	"fmt"
	"strconv"
)

type leaderboardConfig struct {
	name    string
	sqlStmt string
	prefix  string
}

var Config struct {
	mainChannel string
	db          *sql.DB
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

func generateLeaderboardMessage(prefix string, rows *sql.Rows) (string, error) {
	message := prefix
	place := 1
	for rows.Next() {
		var userID string
		var totalPoints int
		var specialPoints int
		err := rows.Scan(&userID, &totalPoints, &specialPoints)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}
		message += fmt.Sprintf("%d. %s: %d (%d)\n", place, userID, totalPoints, specialPoints)
		place++
	}
	return message, nil
}

func SavePoints(userID string, points int) error {
	db := Config.db
	sqlStmt := `
		SELECT COUNT(*) FROM points
		WHERE user_id = ? AND date(timestamp) = date('now');
	`
	var count int
	err := db.QueryRow(sqlStmt, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing points: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("points already recorded today for user %s", userID)
	}

	sqlStmt = `
		INSERT INTO points (timestamp, user_id, points) 
		VALUES (datetime('now'), ?, ?);
	`
	_, err = db.Exec(sqlStmt, userID, points)
	if err != nil {
		return fmt.Errorf("failed to save points: %w", err)
	}

	return nil
}

func formatNumber(n int) string {
	str := strconv.Itoa(n)
	if len(str) < 4 {
		return str
	}

	var result []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
