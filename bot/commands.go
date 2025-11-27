package main

import (
	"log"
	"strings"

	inter "github.com/N1ktarchik/Wishlist_bot/interaction"
	keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandUpdate(update tgbot.Update, bot *tgbot.BotAPI) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover from panic%d", r)
		}
	}()

	if update.CallbackQuery != nil {
		inter.ButtonProcessing(update, bot, *update.CallbackQuery)
		return
	}

	if update.Message != nil && update.Message.Text != "" {

		if update.Message.From == nil {
			return
		}

		user := update.Message.From
		text := update.Message.Text

		switch {

		case text == "/start":
			StartMessage := "Hi, " + user.UserName + " ,I'm waiting for your most explicit desires."
			msg := tgbot.NewMessage(update.Message.Chat.ID, StartMessage)
			bot.Send(msg)

		case strings.HasPrefix(text, "/friend"):

			mas := strings.Split(text, " ")

			if len(mas) != 2 {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "The tag is not recognized! Please send it like this: /friend tag_friend")
				bot.Send(msg)
				keyboard.Menu(update, bot)
				return
			}

			//name:=mas[1]

			//функция поиска друга

		case text == "/menu":
			keyboard.Menu(update, bot)
			return

		case text == "➕ Add new wish":
		//обработка через БД

		case text == "❌ Delete wish":
			//обработка через БД

		case text == "✏️ Change wish":
			//обработка через БД

		case text == "➡️ Next wish":
			//обработка через БД
			//Обдумать архитектуру

		case text == "🔙 Exit to main menu":
			keyboard.Menu(update, bot)
			return

		case text == "✅ Reserve wish":
			//обработка через БД

		default:
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Command not faund!")
			bot.Send(msg)

		}

		keyboard.Menu(update, bot)

	}
}
