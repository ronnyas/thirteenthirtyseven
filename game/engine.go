package game

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/language"
	"github.com/ronnyas/thirteenthirtyseven/leet"
)

var Config struct {
	mainChannel string
	db          *sql.DB
}

func StartEngine(s *discordgo.Session) {
	db := Config.db
	mainChannel := Config.mainChannel
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
				panic(err)
			}
			defer rows.Close()

			leaderboardMessage, err := leet.GenerateLeaderboardMessage(
				language.GetTranslation("game_gen_lb_msg"),
				rows,
			)
			if err != nil {
				log.Fatal(err)
				continue
			}

			s.ChannelMessageSend(mainChannel, leaderboardMessage)

			// update streaks
			_, brokenStreaks, err := leet.UpdateAllStreaks(db)
			if err != nil {
				log.Fatal(err)
				continue
			}
			for _, brokenStreak := range brokenStreaks {
				s.ChannelMessageSend(mainChannel, fmt.Sprintf(language.GetTranslation("game_lb_broke_streak"), brokenStreak.UserID, brokenStreak.Duration()))
			}

			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
