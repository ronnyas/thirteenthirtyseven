package coinflip

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".flip" || m.Content == ".flip mynt" || m.Content == ".flip krone" {
		var outcome string
		flip := rand.Intn(2)

		switch m.Content {
		case ".flip mynt":
			outcome = "Mynt"
			if flip == 0 {
				s.ChannelMessageSend(m.ChannelID, "Gratulerer du vant.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Desverre du tapte.")
			}
		case ".flip krone":
			outcome = "Krone"
			if flip == 0 {
				s.ChannelMessageSend(m.ChannelID, "Desverre du tapte.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Gratulerer du vant.")
			}
		case ".flip":
			if flip == 0 {
				outcome = "Mynt"
			} else {
				outcome = "Krone"
			}
		}

		s.ChannelMessageSend(m.ChannelID, outcome)
	}

}
