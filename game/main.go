package game

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type leaderboardConfig struct {
	name    string
	sqlStmt string
	prefix  string
}

var Config struct {
	mainChannel string
	db 			*sql.DB
}

func SetDatabase(db *sql.DB) {
	Config.db = db
}
func SetMainChannel(channelID string) {
	Config.mainChannel = channelID
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

func SavePoints(userID string, points int) {
	db := Config.db
	sqlStmt := `
		select count(*) from points
		where user_id = ? and timestamp >= date('now');
		`
	var count int
	err := db.QueryRow(sqlStmt, userID).Scan(&count)
	if err != nil {
		panic(err)
	}
	
	if count > 0 {
		return
	} 
		
	sqlStmt = `
	insert into points (timestamp, user_id, points) 
	values (?, ?, ?);
	`
	_, err = db.Exec(sqlStmt, time.Now(), userID, points)
	if err != nil {
		panic(err)
	}
}

func DailyScoreReport(s *discordgo.Session) {
	db := Config.db
	mainChannel := Config.mainChannel
	log.Println("DailyScoreReport started")
	var last_report string = ""
	for {
		current_time := time.Now()
		if current_time.Hour() == 13 && current_time.Minute() == 38 {
			log.Println("DailyScoreReport running")
			// check if it's already been posted
			if last_report == current_time.Format("2006-01-02") {
				continue
			}
			last_report = current_time.Format("2006-01-02")
			log.Println(last_report)
	
			sqlStmt := `
				select user_id, sum(points) from points
				where timestamp >= date('now', 'start of day')
				group by user_id
				order by sum(points) desc
				limit 10;
			`
			rows, err := db.Query(sqlStmt)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			leaderboardMessage, err := generateLeaderboardMessage(
				"Time's up! Here's todays points:\n",
				rows,
			)
			if err != nil {
				log.Fatal(err)
				continue
			}

			s.ChannelMessageSend(mainChannel, leaderboardMessage)

			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "1337" {
		current_time := time.Now()

		if current_time.Hour() != 13 || current_time.Minute() != 37 {
			return
		}
		
		points := calculatePointsFromTimestamp(m.Timestamp)
		
		SavePoints(m.Author.Username, points)
		
		s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
	}

	if m.Content == ".time" {
		s.ChannelMessageSend(m.ChannelID, time.Now().Format("2006-01-02 15:04:05"))
	}

	if m.Content == "1337 lb" {
		db := Config.db

		leaderboardConfigs := []leaderboardConfig{
			{
				name:    "all time",
				sqlStmt: "select user_id, sum(points) from points group by user_id order by sum(points) desc limit 10;",
				prefix:  "\n\n**Leaderboard all time:**\n",
			},
			{
				name:    "this week",
				sqlStmt: "select user_id, sum(points) from points where timestamp >= date('now', '-7 day') group by user_id order by sum(points) desc limit 10;",
				prefix:  "\n\n**Leaderboard this week:**\n",
			},
		}
	
		for _, config := range leaderboardConfigs {
			rows, err := db.Query(config.sqlStmt)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
	
			leaderboardMessage, err := generateLeaderboardMessage(config.prefix, rows)
			if err != nil {
				panic(err)
			}
	
			err = rows.Err()
			if err != nil {
				panic(err)
			}
	
			s.ChannelMessageSend(m.ChannelID, leaderboardMessage)
		}
	}
}