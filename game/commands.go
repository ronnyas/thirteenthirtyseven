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
		output += ".help\n"
		output += ".setup\n"
		output += ".time\n"
		output += "1337\n"
		output += "1337 lb\n"
		output += "1337 streak\n"
		output += "1337 setmain <channelid>\n"
		output += "1337 setactive <true|false>\n"
		output += "1337 setstreak <3>\n"
		output += ".getchannelid\n"
		output += ".getserverid\n"
		output += ".norris\n"
		output += ".flip\n"
		output += ".flip <mynt|krone>\n"

		s.ChannelMessageSend(m.ChannelID, output)
	}

	if m.Content == ".getchannelid" {
		s.ChannelMessageSend(m.ChannelID, m.ChannelID)
	}

	if m.Content == ".getserverid" {
		s.ChannelMessageSend(m.ChannelID, m.GuildID)
	}
}
