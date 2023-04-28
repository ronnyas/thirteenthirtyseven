package coinflip

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/language"
)

/*
func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".test" {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Println("channel error: ", err)
			return
		}

		if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
			s.ChannelMessageSend(m.ChannelID, "Private send")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Channel send")

		}
	}
}*/

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	head := ".flip " + language.GetTranslation("coin_lhead")
	tail := ".flip " + language.GetTranslation("coin_ltail")

	if m.Content == ".flip" || m.Content == head || m.Content == tail {
		var outcome string
		flip := rand.Intn(2)

		switch m.Content {
		case head:
			if flip == 0 {
				s.ChannelMessageSend(m.ChannelID, language.GetTranslation("coin_congrats"))
				outcome = language.GetTranslation("coin_head")
			} else {
				s.ChannelMessageSend(m.ChannelID, language.GetTranslation("coin_sorry"))
				outcome = language.GetTranslation("coin_tail")
			}
		case tail:
			if flip == 0 {
				s.ChannelMessageSend(m.ChannelID, language.GetTranslation("coin_sorry"))
				outcome = language.GetTranslation("coin_head")
			} else {
				s.ChannelMessageSend(m.ChannelID, language.GetTranslation("coin_congrats"))
				outcome = language.GetTranslation("coin_tail")
			}
		case ".flip":
			if flip == 0 {
				outcome = language.GetTranslation("coin_head")
			} else {
				outcome = language.GetTranslation("coin_tail")
			}
		}

		s.ChannelMessageSend(m.ChannelID, outcome)
	}

}
