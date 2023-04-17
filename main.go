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
	"github.com/ronnyas/thirteenthirtyseven/chat"
	"github.com/ronnyas/thirteenthirtyseven/coinflip"
	"github.com/ronnyas/thirteenthirtyseven/config"
	"github.com/ronnyas/thirteenthirtyseven/database"
	"github.com/ronnyas/thirteenthirtyseven/game"
	"github.com/ronnyas/thirteenthirtyseven/leet"
	"github.com/ronnyas/thirteenthirtyseven/norris"
)

func main() {
	log.Println("Loading config")
	cfg := config.LoadConfig()

	log.Println("Starting bot")
	log.Println("Config:")

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
		log.Fatal("Unable to initialize discord session,", err)
		return
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Content == ".time" {
			s.ChannelMessageSend(m.ChannelID, time.Now().Format("2006-01-02 15:04:05"))
		}
	})

	discord.AddHandler(leet.Commands)
	leet.SetDatabase(db)
	leet.SetMainChannel(cfg.MainChannel)
	leet.SetStreakDays(cfg.StreakDays)

	discord.AddHandler(game.Commands)
	game.SetDatabase(db)
	game.SetMainChannel(cfg.MainChannel)

	discord.AddHandler(chat.Commands)
	chat.SetOpenAIKey(cfg.OpenAIKey)

	discord.AddHandler(coinflip.Commands)
	discord.AddHandler(norris.Commands)

	// discord.AddHandler(chat.Commands)
	// chat.SetOpenAIKey(cfg.OpenAIKey)

	// temp code
	// check if there are any data in the streaks table. if not , run BackfillStreaks
	sqlStmt := `select id from streaks limit 1;`
	row := db.QueryRow(sqlStmt)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Println("No streaks found, backfilling")
		leet.BackfillStreaks(db)
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		log.Fatal("Can't connect to discord: ", err)
		return
	}

	go game.StartEngine(discord)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
