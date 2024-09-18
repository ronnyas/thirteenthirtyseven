package game

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/rand"
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
		points := calculatePointsFromTimestamp(m.Timestamp)
		save := SavePoints(m.Author.Username, points)
		if save {
			s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
		}

	}

	// 420 Easter Egg
	if m.Content == "420" {
		current_time := time.Now()

		// 4:20am on any day gives half points
		if current_time.Hour() == 4 && current_time.Minute() == 20 {
			points := (calculatePointsFromTimestamp(m.Timestamp) / 2)

			// If it's also on the 20th of april, give 420 extra points
			if current_time.Month() == time.April && current_time.Day() == 20 {
				points += 420
			}

			save := SavePoints(m.Author.Username, points)
			if save {
				s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
			}
		}
	}

	if m.Content == "1234" {
		current_time := time.Now()

		// 1234 on any day gives 1-4 point for lolz
		if current_time.Hour() == 12 && current_time.Minute() == 34 {
			points := 0
			points += rand.Intn(4) + 1

			save := SavePoints(m.Author.Username, points)
			if save {
				s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
			}
		}
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