package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/N1ktarchik/Wishlist_bot/database"
	inter "github.com/N1ktarchik/Wishlist_bot/interaction"
	keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandUpdate(update tgbot.Update, bot *tgbot.BotAPI, db *sql.DB) {

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
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		status, err := database.GetUserStatusByID(db, chatID)
		if err != nil {
			log.Print(err)
		}

		if status != nil && status.Step != 0 && status.IsAlive() {
			inter.ProcessingNewWish(status, update, bot, db)
			return
		}

		if status != nil && status.Step != 0 && !status.IsAlive() {
			status.Delete(db)
			bot.Send(tgbot.NewMessage(chatID, "⏳ Время добавления истекло. Начните заново."))
		}

		switch {

		case text == "/start":
			StartMessage := "Привет, " + user.UserName + " ,я жду ваших самых откровенных желаний."
			msg := tgbot.NewMessage(chatID, StartMessage)
			bot.Send(msg)

		case strings.HasPrefix(text, "/friend"):

			mas := strings.Split(text, " ")

			if len(mas) != 2 {
				msg := tgbot.NewMessage(chatID, "Тег не распознается! Пожалуйста, отправьте его таким образом: /friend тег_друга")
				bot.Send(msg)
				keyboard.Menu(chatID, bot)
				return
			}

			//name:=mas[1]

			//функция поиска друга

		case text == "/menu":
			keyboard.Menu(chatID, bot)
			return

		case text == "➕ Добавить новое желание":
			//обработка через БД
			user := database.User{ChatID: chatID, UserName: fmt.Sprint(user.UserName)}
			err := user.AddToDB(db)
			if err != nil {
				log.Printf("writing user to DB error. %v", err)
				msg := tgbot.NewMessage(chatID, "Error! Send the screenshot to adminnistrator.")
				bot.Send(msg)
				keyboard.Menu(chatID, bot)
				return
			}

			inter.HandleAddNewWish(chatID, bot, db)
			return

		case text == "❌ Удалить желание":
			//обработка через БД

		case text == "✏️ Изменить желание":
			//обработка через БД

		case text == "➡️ Следующее желание":
			//обработка через БД
			//Обдумать архитектуру

		case text == "⬅️ Предыдущие желание":

		case text == "🔙 Вернуться в главное меню":
			keyboard.Menu(chatID, bot)
			return

		case text == "✅ Зарезервировать желание":
			//обработка через БД

		default:
			msg := tgbot.NewMessage(chatID, "Такой команды не существует!")
			bot.Send(msg)

		}

		keyboard.Menu(chatID, bot)

	}
}
