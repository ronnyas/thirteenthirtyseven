package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bwmarrin/discordgo"
	"github.com/ronnyas/thirteenthirtyseven/chat/norris"
	"github.com/ronnyas/thirteenthirtyseven/config"
	"github.com/ronnyas/thirteenthirtyseven/database"
	"github.com/ronnyas/thirteenthirtyseven/game"
	"github.com/ronnyas/thirteenthirtyseven/game/coinflip"
	"github.com/ronnyas/thirteenthirtyseven/game/leet"
	"github.com/ronnyas/thirteenthirtyseven/language"
)

func main() {

	language.SetLanguage("no")

	log.Println(language.GetTranslation("main_config_load"))
	cfg := config.LoadConfig()

	log.Println(language.GetTranslation("main_starting_bot"))
	log.Println(language.GetTranslation("main_config"))

	key := reflect.ValueOf(cfg).Elem()
	for i := 0; i < key.NumField(); i++ {
		field := key.Field(i)
		log.Println("\t" + key.Type().Field(i).Name + ": " + fmt.Sprintf("%v", field.Interface()))
	}

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.SetupDatabaseSchema(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	discord, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal(language.GetTranslation("main_cant_init_discord"), err)
		return
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Content == ".time" {
			discordTime := m.Timestamp
			systemTime := time.Now()
			diff := systemTime.Sub(discordTime)

			s.ChannelMessageSend(m.ChannelID, "system time: "+systemTime.String())
			s.ChannelMessageSend(m.ChannelID, "discord time: "+discordTime.String())
			s.ChannelMessageSend(m.ChannelID, "time difference: "+diff.String())
		}

		if m.Content == ".setup" {
			server, err := s.Guild(m.GuildID)
			if err != nil {
				log.Println("s.state.guild: ", err)
				return
			}

			if m.Author.ID == server.OwnerID {
				checkCfg := db.QueryRow("SELECT id FROM `config` WHERE `serverid` = " + m.GuildID)
				var id int
				err = checkCfg.Scan(&id)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "No configdata, Inserting...")
					qStmt := "INSERT INTO config (serverid, name, value) VALUES(?, ?, ?)"
					type ConfigStuff struct {
						Name  string
						Value string
					}

					var config []ConfigStuff = []ConfigStuff{
						{
							Name:  "leet_mainchannel",
							Value: "false",
						},
						{
							Name:  "leet_active",
							Value: "false",
						},
						{
							Name:  "leet_streakdays",
							Value: "3",
						},
						{
							Name:  "mafia_mainchannel",
							Value: "false",
						},
						{
							Name:  "mafia_active",
							Value: "false",
						},
					}

					for _, conf := range config {
						_, err := db.Exec(qStmt, m.GuildID, conf.Name, conf.Value)
						if err != nil {
							log.Println("Can't insert default config:", err)
							s.ChannelMessageSend(m.ChannelID, "Something went wrong.")
							return
						}
					}
					s.ChannelMessageSend(m.ChannelID, "Done inserting to the database.")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Only the server owner/creator can run this")
			}
		}
	})

	discord.AddHandler(leet.Commands)
	discord.AddHandler(game.Commands)
	discord.AddHandler(coinflip.Commands)
	discord.AddHandler(norris.Commands)

	leet.SetDatabase(db)
	game.SetDatabase(db)

	// temp code
	// check if there are any data in the streaks table. if not , run BackfillStreaks
	sqlStmt := `select id from streaks limit 1;`
	row := db.QueryRow(sqlStmt)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Println(language.GetTranslation("main_no_streak_backfill"))
		leet.BackfillStreaks(db)
	}

	// https://discord-intents-calculator.vercel.app/
	// guild_presences, guild_messages, guild_message_reactions, direct_messages, direct_message_reactions
	discord.Identify.Intents = 13824

	err = discord.Open()
	if err != nil {
		log.Fatal(language.GetTranslation("main_cant_conn_discord"), err)
		return
	}

	go game.StartEngine(discord)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
