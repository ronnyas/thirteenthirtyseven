package game

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/database"
)

func Commands(s DiscordSession, m *discordgo.MessageCreate) {
	if s.IsBotUser(m.Author.ID) {
		return
	}

	if m.ChannelID != Config.mainChannel {
		return
	}

	currentTime := getCurrentTime()

	switch m.Content {
	case "1337":
		handleRegularPoints(s, m, currentTime)
	case "420":
		handleEasterEgg(s, m, currentTime)
	case "1234":
		handleStaticPoints(s, m, currentTime)
	case "1337 lb":
		handleLeaderboard(s, m)
	case "1337 streak":
		handleStreaks(s, m)
	}
}

func handleRegularPoints(s DiscordSession, m *discordgo.MessageCreate, currentTime time.Time) {
	if !IsValid1337Time(currentTime) {
		return
	}

	points := Calculate1337Points(currentTime)
	err := database.AddPoints(Config.db, m.Author.Username, currentTime, points)
	if err == nil {
		addReaction(s, m)
	} else {
		log.Printf("Failed to add points for %s: %v", m.Author.Username, err)
	}
}

func handleEasterEgg(s DiscordSession, m *discordgo.MessageCreate, currentTime time.Time) {
	if !IsValidSpecialTime(currentTime) {
		return
	}

	points := CalculateSpecialPoints(currentTime)
	eggType := GetEasterEggType(currentTime)
	err := database.AddSpecialPoints(Config.db, m.Author.Username, currentTime, points, eggType)
	if err == nil {
		addReaction(s, m)
	} else {
		log.Printf("Failed to add special points for %s: %v", m.Author.Username, err)
	}
}

func handleStaticPoints(s DiscordSession, m *discordgo.MessageCreate, currentTime time.Time) {
	err := database.AddSpecialPoints(Config.db, m.Author.Username, currentTime, 10, "1234")
	if err == nil {
		addReaction(s, m)
	} else {
		log.Printf("Failed to add static points for %s: %v", m.Author.Username, err)
	}
}

func handleLeaderboard(s DiscordSession, m *discordgo.MessageCreate) {
	leaderboardConfigs := []leaderboardConfig{
		{
			name: "all time",
			sqlStmt: `
				WITH combined AS (
					SELECT user_id, 
						SUM(CASE WHEN type = 'regular' THEN points ELSE 0 END) as regular_points,
						SUM(CASE WHEN type = 'special' THEN points ELSE 0 END) as special_points
					FROM (
						SELECT user_id, points, 'regular' as type FROM points
						UNION ALL
						SELECT user_id, points, 'special' as type FROM special_points
					)
					GROUP BY user_id
				)
				SELECT user_id, 
					(regular_points + special_points) as total_points,
					special_points
				FROM combined
				ORDER BY total_points DESC
				LIMIT 10;
			`,
			prefix: "\n\n**Leaderboard all time:**\n",
		},
		{
			name: "this week",
			sqlStmt: `
				WITH combined AS (
					SELECT user_id, 
						SUM(CASE WHEN type = 'regular' THEN points ELSE 0 END) as regular_points,
						SUM(CASE WHEN type = 'special' THEN points ELSE 0 END) as special_points
					FROM (
						SELECT user_id, points, 'regular' as type 
						FROM points 
						WHERE date(timestamp) >= date('now', 'weekday 0', '-6 days')
						UNION ALL
						SELECT user_id, points, 'special' as type 
						FROM special_points 
						WHERE date(timestamp) >= date('now', 'weekday 0', '-6 days')
					)
					GROUP BY user_id
				)
				SELECT user_id, 
					(regular_points + special_points) as total_points,
					special_points
				FROM combined
				ORDER BY total_points DESC
				LIMIT 10;
			`,
			prefix: "\n\n**Leaderboard this week:**\n",
		},
	}

	for _, config := range leaderboardConfigs {
		rows, err := Config.db.Query(config.sqlStmt)
		if err != nil {
			log.Printf("Error querying leaderboard: %v", err)
			continue
		}
		defer rows.Close()

		leaderboardMessage, err := generateLeaderboardMessage(config.prefix, rows)
		if err != nil {
			log.Printf("Error generating leaderboard message: %v", err)
			continue
		}

		if len(leaderboardMessage) == len(config.prefix) {
			leaderboardMessage += "No points yet!"
		}

		s.ChannelMessageSend(m.ChannelID, leaderboardMessage)
	}
}

func handleStreaks(s DiscordSession, m *discordgo.MessageCreate) {
	streaks, err := GetActiveStreaks(Config.db)
	if err != nil {
		log.Printf("Error getting active streaks: %v", err)
		return
	}

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

func addReaction(s DiscordSession, m *discordgo.MessageCreate) {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "1337:1079824982613442580")
	if err != nil {
		log.Printf("Failed to add reaction: %v", err)
	}
}
