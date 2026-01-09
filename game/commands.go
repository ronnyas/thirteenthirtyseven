package game

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != Config.mainChannel {
		return
	}

	// Check if message is a valid 1337 variant (case-insensitive)
	validVariants := []string{
		"1337",
		"trettentrettisju",
		"13:37",
		"thirteenthirtyseven",
		"leet",
		"elite",
	}
	
	isValidVariant := false
	contentLower := strings.ToLower(strings.TrimSpace(m.Content))
	for _, variant := range validVariants {
		if strings.ToLower(variant) == contentLower {
			isValidVariant = true
			break
		}
	}
	
	if isValidVariant {
		current_time := time.Now()

		if current_time.Hour() != 13 || current_time.Minute() != 37 {
			return
		}
		
		points := calculatePointsFromTimestamp(m.Timestamp)
		
		save := SavePoints(m.Author.Username, points)
		if save {
			s.MessageReactionAdd(m.ChannelID, m.ID, Config.ReactEmoji)
		}

	}

	if strings.HasPrefix(m.Content, "1337 lb") {
		log.Printf("Leaderboard command received: %q", m.Content)
		db := Config.db
		
		// Parse the command to extract year if provided
		parts := strings.Fields(m.Content)
		var year int
		var yearStr string
		var showTotal bool
		
		if len(parts) == 3 {
			// Check if it's "total" or a year
			yearStr = parts[2]
			if strings.ToLower(yearStr) == "total" {
				showTotal = true
			} else {
				// Year provided: "1337 lb <year>"
				parsedYear, err := strconv.Atoi(yearStr)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid year: %s. Please provide a valid year (e.g., 2024) or 'total'.", yearStr))
					return
				}
				year = parsedYear
			}
		} else if len(parts) == 2 {
			// No year provided: "1337 lb" - use current year
			year = time.Now().Year()
		} else {
			// Invalid command format
			s.ChannelMessageSend(m.ChannelID, "Usage: `1337 lb`, `1337 lb <year>`, or `1337 lb total`")
			return
		}
		
		var sqlStmt string
		var prefix string
		
		if showTotal {
			// All-time leaderboard
			sqlStmt = "select user_id, sum(points) from points group by user_id order by sum(points) desc limit 10;"
			prefix = "**Leaderboard all time:**\n"
		} else {
			// Validate year range (reasonable bounds)
			if year < 2000 || year > 2100 {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid year: %d. Year must be between 2000 and 2100.", year))
				return
			}
			
			// Build SQL query to filter by year
			sqlStmt = fmt.Sprintf(
				"select user_id, sum(points) from points where strftime('%%Y', timestamp) = '%d' group by user_id order by sum(points) desc limit 10;",
				year,
			)
			prefix = fmt.Sprintf("**Leaderboard %d:**\n", year)
		}
		
		rows, err := db.Query(sqlStmt)
		if err != nil {
			if showTotal {
				log.Printf("Error querying database for total leaderboard: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Error querying all-time leaderboard.")
			} else {
				log.Printf("Error querying database for year %d: %v", year, err)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error querying leaderboard for year %d.", year))
			}
			return
		}
		defer rows.Close()

		leaderboardMessage, err := generateLeaderboardMessage(prefix, rows)
		if err != nil {
			if showTotal {
				log.Printf("Error generating leaderboard message: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Error generating all-time leaderboard.")
			} else {
				log.Printf("Error generating leaderboard message: %v", err)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error generating leaderboard for year %d.", year))
			}
			return
		}

		err = rows.Err()
		if err != nil {
			if showTotal {
				log.Printf("Error iterating rows: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Error processing all-time leaderboard.")
			} else {
				log.Printf("Error iterating rows: %v", err)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error processing leaderboard for year %d.", year))
			}
			return
		}

		if len(leaderboardMessage) == len(prefix) {
			leaderboardMessage += "No points yet!"
		}

		s.ChannelMessageSend(m.ChannelID, leaderboardMessage)
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