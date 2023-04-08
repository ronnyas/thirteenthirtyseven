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
			
			// update streaks
			_, brokenStreaks, err := UpdateAllStreaks(db)
			if err != nil {
				log.Fatal(err)
				continue
			}
			for _, brokenStreak := range brokenStreaks {
				s.ChannelMessageSend(mainChannel, fmt.Sprintf("%s broke their streak of %d days", brokenStreak.UserID, brokenStreak.Duration()))
			}

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

	if m.ChannelID != Config.mainChannel {
		return
	}


	if m.Content == "1337" {
		current_time := time.Now()

		if current_time.Hour() != 13 || current_time.Minute() != 37 {
			return
		}
		
		points := calculatePointsFromTimestamp(m.Timestamp)
		
		save := SavePoints(m.Author.Username, points)
		if save {
			s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
		}

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
				sqlStmt: "select user_id, sum(points) from points where date(timestamp) >= date('now', 'weekday 0', '-6 days') group by user_id order by sum(points) desc limit 10;",
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

			if len(leaderboardMessage) == len(config.prefix) {
				leaderboardMessage += "No points yet!"
			}

	
			s.ChannelMessageSend(m.ChannelID, leaderboardMessage)
		}
	}

	if m.Content == "1337 streak" {
		streaks, err := GetActiveStreaks(Config.db)
		if err != nil {
			log.Fatal(err)
			return
		}
		// check if there are any active streaks
		if len(streaks) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No active streaks :(")
			return
		}
		streakMsg := "Active streaks:\n"
		for _, streak := range streaks {
			streakDuration := streak.Duration()
			streakMsg += fmt.Sprintf("%s: %d days\n", streak.UserID, streakDuration)
		}

		s.ChannelMessageSend(m.ChannelID, streakMsg)
	
	}
}