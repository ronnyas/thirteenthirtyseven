package coinflip

var Config struct {
	mainChannel string
}

func SetMainChannel(channelID string) {
	Config.mainChannel = channelID
}
