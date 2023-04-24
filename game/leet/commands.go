package leet

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/config"
	"github.com/ronnyas/thirteenthirtyseven/database"
	"github.com/ronnyas/thirteenthirtyseven/language"
)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server, err := s.Guild(m.GuildID)
	if err != nil {
		log.Println("s.state.guild: ", err)
		return
	}

	args := strings.Split(m.Content, " ")
	var command string
	if len(command) > 2 {
		command = args[0] + " " + args[1]
	}

	if (len(args)) == 3 {
		if m.Author.ID == server.OwnerID {

			query := `UPDATE "config" SET "value" = ? WHERE serverid = ? AND name = ?`

			if command == "1337 setmain" {
				db.Exec(query, args[2], m.GuildID, "leet_mainchannel")
				//SetMainChannel(args[2])
				s.ChannelMessageSend(m.ChannelID, "Main channel set!")
			}

			if command == "1337 setactive" {
				if args[2] == "true" {
					db.Exec(query, args[2], m.GuildID, "leet_active")
					//SetActive(args[2])
					s.ChannelMessageSend(m.ChannelID, "1337 Activated")
				} else if args[2] == "false" {
					db.Exec(query, args[2], m.GuildID, "leet_active")
					//SetActive(args[2])
					s.ChannelMessageSend(m.ChannelID, "1337 Deactivated")
				}
			}

			if command == "1337 setstreak" {
				setStreak, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatal(err)
					return
				}
				_, err = db.Exec(query, setStreak, m.GuildID, "leet_streakdays")
				if err != nil {
					log.Println(err)
					return
				}
				//SetStreakDays(setStreak)
				s.ChannelMessageSend(m.ChannelID, "1337 Streak set to "+args[2])
			}
		}
	}

	mainChannel, err := GetMainChannel(m.GuildID)
	if err != nil {
		log.Println("Error getting mainchannel to checkmainchannel: ", err)
	}
	if mainChannel != m.ChannelID {
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
