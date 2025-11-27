package main

import (
	"errors"
	"log"
	"os"

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

	for update := range updates {

		go CommandUpdate(update, bot)
	}
}

func GetToken(BOT_TOKEN string) (*tgbot.BotAPI, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("Error loading .env file")
	}
	bot, err := tgbot.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, errors.New("Token reading error")
	}

	return bot, nil
}
