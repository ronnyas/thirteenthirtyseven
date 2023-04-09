package chat

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	Timestamp time.Time
	Username string
	Message string
}

var Messages []Message
var lastMesssageSent = time.Now().Add(-3 * time.Hour)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != Config.mainChannel {
		return
	}

	Messages = RefreshChatlog(
		Message{
			Timestamp: m.Timestamp,
			Username: m.Author.Username,
			Message: m.Content,
		},
	)

	// send no more than one message every 3 hours
	if time.Now().Sub(lastMesssageSent).Hours() < 3 {
		return
	}

	if len(Messages) > 3 {
		totalLength := 0
		for _, message := range Messages {
			totalLength += len(message.Message)
		}

		if totalLength > 150 {
			answerFromOpenAI := GenerateAnswerToDiscussion(Messages)

			s.ChannelMessageSend(m.ChannelID, answerFromOpenAI)

			Messages = []Message{}
			lastMesssageSent = time.Now()
		}

	}
}

func RefreshChatlog(m Message) []Message {
	if len(m.Message) > 8 && m.Message[:8] == "https://" {
		return Messages
	}

	Messages = append(Messages, m)

	if len(Messages) > 1 {
		cutoffTime := time.Now().Add(-30 * time.Minute)

		// remove messages older than 30 minutes
		for i, message := range Messages {
			if message.Timestamp.Before(cutoffTime) {
				Messages = append(Messages[:i], Messages[i+1:]...)
			}
		}
	}

	return Messages
}