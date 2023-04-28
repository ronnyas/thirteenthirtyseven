package chat

import (
	"context"
	"database/sql"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var Config struct {
	mainChannel string
	db 			*sql.DB
	OpenAIKey  	string
}

func SetDatabase(db *sql.DB) {
	Config.db = db
}
func SetMainChannel(channelID string) {
	Config.mainChannel = channelID
}

func SetOpenAIKey(key string) {
	Config.OpenAIKey = key
}


func GenerateAnswerToDiscussion(messages []Message) string {
	prompt := `Du er en av personene i samtalen og du skal være sarkastisk med dine svar. Svar med 1-3 setninger.  Chatlog:`
	for _, message := range messages {
		prompt += message.Username + ": " + message.Message + "\n"
	}

	client := openai.NewClient(Config.OpenAIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}
