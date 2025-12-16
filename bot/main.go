package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/N1ktarchik/Wishlist_bot/database"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	bot, err := GetToken(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("connect to DB error. %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("connect to DB error. %v", err)
		}
	}()

	err = database.CreateTables(db)
	if err != nil {
		log.Fatalf("creating tables in DB error. %v", err)
	}

	go CleanUpService(db)

	for update := range updates {

		go CommandUpdate(update, bot, db)
	}
}

func GetToken(BOT_TOKEN string) (*tgbot.BotAPI, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading .env file")
	}
	bot, err := tgbot.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, errors.New("token reading error")
	}

	return bot, nil
}

func CleanUpService(db *sql.DB) {

	err := database.CleanOverdueStatuses(db)
	if err != nil {
		log.Printf("error clean service. %v", err)
	}

	err = database.CleanExpiredSessions(db)
	if err != nil {
		log.Printf("error clean service. %v", err)
	}

	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		err := database.CleanOverdueStatuses(db)
		if err != nil {
			log.Printf("error clean service. %v", err)
		}

		err = database.CleanExpiredSessions(db)
		if err != nil {
			log.Printf("error clean service. %v", err)
		}
	}
}
