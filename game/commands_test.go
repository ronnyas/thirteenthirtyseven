package game

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/database"
)

type mockSession struct {
	lastReaction    string
	lastMessageSent string
	lastChannelID   string
	reactionAdded   bool
	messageSent     bool
	botID           string
}

func (m *mockSession) MessageReactionAdd(channelID, messageID, emoji string) error {
	m.lastReaction = emoji
	m.lastChannelID = channelID
	m.reactionAdded = true
	return nil
}

func (m *mockSession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
	m.lastMessageSent = content
	m.lastChannelID = channelID
	m.messageSent = true
	return nil, nil
}

func (m *mockSession) IsBotUser(userID string) bool {
	return userID == m.botID
}

func TestCommandHandling(t *testing.T) {
	// Setup test database
	db, err := database.Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := database.SetupDatabaseSchema(db); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	Config.db = db
	Config.mainChannel = "test-channel"

	testCases := []struct {
		name           string
		message        string
		time           time.Time
		expectReaction bool
		expectMessage  bool
		checkPoints    bool
		expectedPoints int
	}{
		{
			name:           "Valid 1337",
			message:        "1337",
			time:           time.Date(2024, 1, 1, 13, 37, 0, 0, time.UTC),
			expectReaction: true,
			checkPoints:    true,
			expectedPoints: 60,
		},
		{
			name:           "Invalid 1337 time",
			message:        "1337",
			time:           time.Date(2024, 1, 1, 13, 38, 0, 0, time.UTC),
			expectReaction: false,
			checkPoints:    true,
			expectedPoints: 0,
		},
		{
			name:           "Valid 420",
			message:        "420",
			time:           time.Date(2024, 1, 1, 16, 20, 0, 0, time.UTC),
			expectReaction: true,
			checkPoints:    true,
			expectedPoints: 6,
		},
		{
			name:           "Special 420",
			message:        "420",
			time:           time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC),
			expectReaction: true,
			checkPoints:    true,
			expectedPoints: 600,
		},
		{
			name:           "Valid 1234",
			message:        "1234",
			time:           time.Date(2024, 1, 1, 12, 34, 0, 0, time.UTC),
			expectReaction: true,
			checkPoints:    true,
			expectedPoints: 10,
		},
		{
			name:           "Leaderboard request",
			message:        "1337 lb",
			time:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expectReaction: false,
			expectMessage:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockSession{}
			msg := &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Content:   tc.message,
					ChannelID: Config.mainChannel,
					Author: &discordgo.User{
						ID:       "test-user",
						Username: "test-user",
					},
				},
			}

			// Set test time
			testTime = &tc.time // You'll need to add this variable to your package

			// Process command
			Commands(mock, msg)

			// Check reaction
			if tc.expectReaction && !mock.reactionAdded {
				t.Error("Expected reaction, but none was added")
			}
			if !tc.expectReaction && mock.reactionAdded {
				t.Error("Did not expect reaction, but one was added")
			}

			// Check message
			if tc.expectMessage && !mock.messageSent {
				t.Error("Expected message, but none was sent")
			}
			if !tc.expectMessage && mock.messageSent {
				t.Error("Did not expect message, but one was sent")
			}

			// Check points if needed
			if tc.checkPoints {
				points, err := database.GetTotalPoints(db, "test-user")
				if err != nil {
					t.Errorf("Failed to get points: %v", err)
				}
				if tc.expectedPoints != -1 && points != tc.expectedPoints {
					t.Errorf("Expected %d points, got %d", tc.expectedPoints, points)
				}
				if tc.expectedPoints == -1 && (points < 1 || points > 4) {
					t.Errorf("Expected points between 1-4, got %d", points)
				}
			}
		})

		// Clear database between tests
		db.Exec("DELETE FROM points")
		db.Exec("DELETE FROM special_points")
	}
}
