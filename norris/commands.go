package norris

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".norris" {
		facts := []string{
			"Chuck Norris doesn't read books. He stares them down until he gets the information he wants.",
			"Chuck Norris can divide by zero.",
			"Chuck Norris can sneeze with his eyes open.",
			"Chuck Norris can speak Braille.",
			"Chuck Norris doesn't do push-ups. He pushes the earth down.",
			"Chuck Norris can make a fire by rubbing two ice-cubes together.",
			"Chuck Norris once won a staring contest against his own reflection.",
			"Chuck Norris can unscramble an egg.",
			"Chuck Norris can kill two stones with one bird.",
			"Chuck Norris can delete the Recycling Bin.",
			"Chuck Norris can slam a revolving door.",
			"Chuck Norris can strangle you with a cordless phone.",
			"When Chuck Norris falls in water, Chuck Norris doesn't get wet. Water gets Chuck Norris.",
			"Chuck Norris once went skydiving, but promised never to do it again. One Grand Canyon is enough.",
			"Chuck Norris can kill your imaginary friends.",
			"Chuck Norris once won a game of Connect Four in three moves.",
			"Chuck Norris can find the needle in the haystack.",
			"Chuck Norris can hear sign language.",
			"Chuck Norris can make onions cry.",
			"Chuck Norris doesn't wear a watch. HE decides what time it is.",
			"Chuck Norris has already been to Mars; that's why there are no signs of life there.",
			"Chuck Norris doesn't mow his lawn, he stands on the porch and dares it to grow.",
			"Chuck Norris once visited the Virgin Islands. They are now The Islands.",
			"Chuck Norris once punched a man in the soul.",
			"If you spell Chuck Norris in Scrabble, you win. Forever.",
			"Chuck Norris breathes air ... five times a day.",
			"In the Beginning there was nothing ... then Chuck Norris roundhouse kicked nothing and told it to get a job.",
			"When God said, “Let there be light!” Chuck Norris said, “Say Please.”",
			"If Chuck Norris were to travel to an alternate dimension in which there was another Chuck Norris and they both fought, they would both win.",
			"The dinosaurs looked at Chuck Norris the wrong way once. You know what happened to them.",
			"Chuck Norris' tears cure cancer. Too bad he has never cried.",
			"Chuck Norris once roundhouse kicked someone so hard that his foot broke the speed of light",
			"Chuck Norris does not sleep. He waits.",
			"Chuck Norris drinks napalm to fight his heartburn.",
			"Chuck Norris' roundhouse kick is so powerful, it can be seen from outer space by the naked eye.",
			"If you want a list of Chuck Norris' enemies, just check the extinct species list.",
			"Chuck Norris once shot an enemy plane down with his finger, by yelling, “Bang!”",
			"Some kids pee their name in the snow. Chuck Norris can pee his name into concrete.",
			"Chuck Norris counted to infinity... twice.",
			"Chuck Norris can do a wheelie on a unicycle.",
			"Once a cobra bit Chuck Norris' leg. After five days of excruciating pain, the cobra died.",
			"The dark is afraid of Chuck Norris.",
			"Chuck Norris can build a snowman out of rain.",
			"Chuck Norris can drown a fish.",
			"When Chuck Norris enters a room, he doesn't turn the lights on, he turns the dark off.",
			"Chuck Norris is the only person that can punch a cyclops between the eye.",
			"Chuck Norris used to beat up his shadow because it was following to close. It now stands 15 feet behind him.",
			"Outer space exists because it's afraid to be on the same planet with Chuck Norris.",
			"Chuck Norris can get in a bucket and lift it up with himself in it.",
			"When Chuck Norris does division, there are no remainders.",
			"Chuck Norris had to stop washing his clothes in the ocean. Too many tsunamis.",
			"Chuck Norris beat the sun in a staring contest.",
			"Chuck Norris can clap with one hand.",
			"Chuck Norris doesn't need to shave. His beard is scared to grow.",
		}

		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(facts))

		s.ChannelMessageSend(m.ChannelID, facts[randomIndex])
	}
}
