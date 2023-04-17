package game

import (
	"github.com/bwmarrin/discordgo"
)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".help" {
		var output string
		output += ".help: Viser denne beskjeden.\n"
		output += "1337: Gir deg et visst antall poeng kl 13:37, baser på sekunder.\n"
		output += "1337 lb: Gir deg en liste over top 10 brukere som har deltatt.\n"
		output += "1337 streak: Gir deg en liste over top 10 brukere som gjør dette hver dag.\n"
		output += ".norris: Gir deg en Chuck Norris fakta.\n"
		output += ".flip: Gir deg mynt eller krone tilbake.\n"
		output += ".flip <mynt|krone>: Spiller mynt eller krone med meg."

		s.ChannelMessageSend(m.ChannelID, output)
	}
}
