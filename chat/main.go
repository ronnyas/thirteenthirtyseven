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
	prompt := `Du er en av personene i samtalen. Svar med 1-3 setninger. Referer til personene i samtalen med navn hvis du vil argumentere for eller imot noe de sier. Kom med nye perspektiver og ikke gjenta det som allerede er sagt.  Vær uenig hvis noe ulogisk er sagt. Engajert og humoristisk uenig eventuelt. Hvis noe er humoristisk eller ironisk, vis samme humor eller ironi. Bruk opp til to emojis. Du kan flette inn et klokt spørsmål om relevant. Unngå generelle spørsmål. Spørsmål som skaper refleksjon=OK. Spørsmål for nysgjerrighetens skyld = Ikke ok. Chatlog:`
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