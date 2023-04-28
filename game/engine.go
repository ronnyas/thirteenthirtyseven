package game

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/game/leet"
	"github.com/ronnyas/thirteenthirtyseven/language"
)

var Config struct {
	mainChannel string
	db          *sql.DB
}

func StartEngine(s *discordgo.Session) {
	db := Config.db

	log.Println(language.GetTranslation("game_started"))
	var last_report string = ""
	for {
		current_time := time.Now()
		if current_time.Hour() == 13 && current_time.Minute() == 38 {
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
				log.Println(err)
			}
			defer rows.Close()

			leaderboardMessage, err := leet.GenerateLeaderboardMessage(
				language.GetTranslation("game_gen_lb_msg"),
				rows,
			)
			if err != nil {
				log.Println(err)
				continue
			}

			servers, err := GetServers()
			if err != nil {
				log.Println("Error getting servers: ", err)
				return
			}

			for _, channels := range servers {
				mainChannel, err := leet.GetMainChannel(channels)
				if err != nil {
					log.Println("Error getting mainchannel to leaderboardmessage: ", err)
				}
				s.ChannelMessageSend(mainChannel, leaderboardMessage)
			}

			// update streaks
			_, brokenStreaks, err := leet.UpdateAllStreaks(db)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, brokenStreak := range brokenStreaks {
				for _, channels := range servers {
					mainChannel, err := leet.GetMainChannel(channels)
					if err != nil {
						log.Println("Error getting mainchannel to brokenstreak: ", err)
						return
					}
					s.ChannelMessageSend(mainChannel, fmt.Sprintf(language.GetTranslation("game_lb_broke_streak"), brokenStreak.UserID, brokenStreak.Duration()))
				}
			}

			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
