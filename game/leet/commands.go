package leet

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/language"
)

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

		points := CalculatePointsFromTimestamp(m.Timestamp)

		save := SavePoints(m.Author.Username, points)
		if save {
			s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
		}
	}

	if m.Content == "1337 lb" {
		db := Config.db

		leaderboardConfigs := []leaderboardConfig{
			{
				name:    "all time",
				sqlStmt: "select user_id, sum(points) from points group by user_id order by sum(points) desc limit 10;",
				prefix:  language.GetTranslation("leet_lb_alltime"),
			},
			{
				name:    "this week",
				sqlStmt: "select user_id, sum(points) from points where date(timestamp) >= date('now', 'weekday 0', '-6 days') group by user_id order by sum(points) desc limit 10;",
				prefix:  language.GetTranslation("leet_lb_week"),
			},
		}

		for _, config := range leaderboardConfigs {
			rows, err := db.Query(config.sqlStmt)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			leaderboardMessage, err := GenerateLeaderboardMessage(config.prefix, rows)
			if err != nil {
				panic(err)
			}

			err = rows.Err()
			if err != nil {
				panic(err)
			}

			if len(leaderboardMessage) == len(config.prefix) {
				leaderboardMessage += language.GetTranslation("leet_lb_nopoints")
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
			s.ChannelMessageSend(m.ChannelID, language.GetTranslation("leet_lb_nostreak"))
			return
		}
		streakMsg := "Active streaks:\n"
		for _, streak := range streaks {
			streakDuration := streak.Duration()
			streakMsg += fmt.Sprintf(language.GetTranslation("leet_days"), streak.UserID, streakDuration)
		}

		s.ChannelMessageSend(m.ChannelID, streakMsg)
	}
}
