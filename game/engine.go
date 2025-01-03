package game

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// For testing
var testTime *time.Time

func getCurrentTime() time.Time {
	if testTime != nil {
		return *testTime
	}
	return time.Now()
}

// DiscordSession interface for mocking in tests
type DiscordSession interface {
	MessageReactionAdd(channelID, messageID, emoji string) error
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
	IsBotUser(userID string) bool
}

type EasterEgg struct {
	Type            string
	ValidateTime    func(time.Time) bool
	CalculatePoints func(time.Time) int
}

var EasterEggs = map[string]EasterEgg{
	"420": {
		Type: "420",
		ValidateTime: func(t time.Time) bool {
			return t.Hour() == 16 && t.Minute() == 20
		},
		CalculatePoints: func(t time.Time) int {
			if t.Hour() == 16 && t.Minute() == 20 {
				return int(float64(Calculate1337Points(time.Date(t.Year(), t.Month(), t.Day(), 13, 37, 0, 0, t.Location()))) * 0.1)
			}
			return 0
		},
	},
	"420_special": {
		Type: "420_special",
		ValidateTime: func(t time.Time) bool {
			return t.Hour() == 16 && t.Minute() == 20 && t.Month() == time.April && t.Day() == 20
		},
		CalculatePoints: func(t time.Time) int {
			if t.Hour() == 16 && t.Minute() == 20 && t.Month() == time.April && t.Day() == 20 {
				return Calculate1337Points(time.Date(t.Year(), t.Month(), t.Day(), 13, 37, 0, 0, t.Location())) * 10
			}
			return 0
		},
	},
}

func IsValid1337Time(t time.Time) bool {
	return t.Hour() == 13 && t.Minute() == 37
}

func Calculate1337Points(t time.Time) int {
	if !IsValid1337Time(t) {
		return 0
	}

	// Start with 60 points at 13:37:00
	// Subtract one point for each second passed
	secondsPassed := t.Second()
	points := 60 - secondsPassed

	// Ensure points don't go negative
	if points < 0 {
		return 0
	}
	return points
}

func IsValidSpecialTime(t time.Time) bool {
	// Check special cases first
	if EasterEggs["420_special"].ValidateTime(t) {
		return true
	}
	// Then check regular cases
	if EasterEggs["420"].ValidateTime(t) {
		return true
	}
	return false
}

func CalculateSpecialPoints(t time.Time) int {
	// Check special cases first
	if EasterEggs["420_special"].ValidateTime(t) {
		return EasterEggs["420_special"].CalculatePoints(t)
	}
	// Then check regular cases
	if EasterEggs["420"].ValidateTime(t) {
		return EasterEggs["420"].CalculatePoints(t)
	}
	return 0
}

func GetEasterEggType(t time.Time) string {
	// Check special cases first
	if EasterEggs["420_special"].ValidateTime(t) {
		return EasterEggs["420_special"].Type
	}
	// Then check regular cases
	if EasterEggs["420"].ValidateTime(t) {
		return EasterEggs["420"].Type
	}
	return ""
}

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
