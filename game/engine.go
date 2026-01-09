package game

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func StartEngine(s *discordgo.Session) {
	db := Config.db
	mainChannel := Config.mainChannel
	log.Println("Game engine started")
	var last_report string = ""
	for {
		current_time := time.Now()
		if current_time.Hour() == 13 && current_time.Minute() == 38 {
			if last_report == current_time.Format("2006-01-02") {
				continue
			}
			last_report = current_time.Format("2006-01-02")
			log.Println(last_report)
	
			// First check if there are any points for today
			checkStmt := `
				select count(*) from points
				where timestamp >= date('now', 'start of day');
			`
			var count int
			err := db.QueryRow(checkStmt).Scan(&count)
			if err != nil {
				log.Printf("Error checking for today's points: %v", err)
				time.Sleep(60 * time.Second)
				continue
			}

			// Only send leaderboard if there were points today
			if count == 0 {
				continue
			}

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