package leet

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
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
	Active      string
}

func VarDump(whatever interface{}) {
	fmt.Printf("%#v\n", whatever)
}

func SetDatabase(db *sql.DB) {
	Config.db = db
}

func GetServers() ([]string, error) {
	query := `SELECT DISTINCT serverid FROM config`
	rows, err := Config.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %s", err)
	}
	defer rows.Close()

	var serverIDs []string
	for rows.Next() {
		var serverID string
		if err := rows.Scan(&serverID); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		serverIDs = append(serverIDs, serverID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}
	return serverIDs, nil
}

func GetMainChannel(serverID string) (string, error) {
	query := `SELECT value FROM config WHERE serverid = ? AND name = 'leet_mainchannel'`
	var mainChannelValue string
	err := Config.db.QueryRow(query, serverID).Scan(&mainChannelValue)
	if err != nil {
		return "", err
	}
	return mainChannelValue, nil
}

func GetStreakDays(serverID string) (string, error) {
	query := `SELECT value FROM config WHERE server = ? AND name = 'leet_streakdays'`
	var streakDaysValue string
	err := Config.db.QueryRow(query, serverID).Scan(&streakDaysValue)
	if err != nil {
		return "", err
	}
	return streakDaysValue, nil
}

func GetServerStatus(serverID string) (string, error) {
	query := `SELECT value FROM config WHERE server = ? AND name = 'leet_active'`
	var serverStatus string
	err := Config.db.QueryRow(query, serverID).Scan(&serverStatus)
	if err != nil {
		return "", err
	}
	return serverStatus, err
}

func CalculatePointsFromTimestamp(timestamp time.Time) int {
	layout := "2006-01-02 15:04:05.000 -0700 MST"

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatalln(err)
	}

	timestamp, err = time.ParseInLocation(layout, timestamp.Format(layout), loc)
	if err != nil {
		log.Fatalln(err)
	}

	points := 60 - timestamp.Second()

	return points
}

func GenerateLeaderboardMessage(prefix string, rows *sql.Rows) (string, error) {
	leaderboardMessage := prefix
	rank := 1

	for rows.Next() {
		var userID string
		var points int

		if err := rows.Scan(&userID, &points); err != nil {
			return "", err
		}

		pointsFormatted := formatNumber(points)

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

		leaderboardMessage += fmt.Sprintf("%s %s\n", userID, pointsFormatted)
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

func formatNumber(number int) string {
	var formattedInt string
	for i, r := range strconv.Itoa(number) {
		if i > 0 && (len(strconv.Itoa(number))-i)%3 == 0 {
			formattedInt += ","
		}
		formattedInt += string(r)
	}

	return formattedInt
}
