package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/N1ktarchik/Wishlist_bot/database"
	inter "github.com/N1ktarchik/Wishlist_bot/interaction"
	"github.com/N1ktarchik/Wishlist_bot/keyboards"
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
		inter.ButtonProcessing(update, bot, *update.CallbackQuery, db)
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
			//принят update/new
			inter.ProcessingNewWish(status, update, bot, db)
			return
		}

		if status != nil && status.Step != 0 && !status.IsAlive() {
			status.Delete(db)
			bot.Send(tgbot.NewMessage(chatID, "⏳ Время добавления истекло. Начните заново."))
		}

		switch {

		case text == "/start":

			user := database.User{ChatID: chatID, UserName: fmt.Sprint(user.UserName)}
			err := user.AddToDB(db)
			if err != nil {
				log.Printf("writing user to DB error. %v", err)
				msg := tgbot.NewMessage(chatID, "Error to add new wish! Send the screenshot to adminnistrator.")
				bot.Send(msg)
				keyboard.Menu(chatID, bot)
				return
			}

			StartMessage := "Привет, " + user.UserName + " ,я жду ваших самых откровенных желаний."
			msg := tgbot.NewMessage(chatID, StartMessage)
			bot.Send(msg)

		case strings.HasPrefix(text, "/friend"):

			mas := strings.Split(text, " ")

			if len(mas) != 2 {
				msg := tgbot.NewMessage(chatID, "Тег не распознается! Пожалуйста, отправьте его таким образом: /friend тег_друга")
				bot.Send(msg)
				return
			}

			friendID, err := database.GetIdByUsername(mas[1], db)
			if err != nil {
				msg := tgbot.NewMessage(chatID, fmt.Sprintf("Пользователь с username %s не найден!", mas[1]))
				log.Print(err)
				bot.Send(msg)
				return
			}

			// session, err := database.GetWishSessonByID(chatID, db)
			// if err != nil {
			// 	log.Print(err)
			// 	return
			// }

			if friendID == chatID {
				msg := tgbot.NewMessage(chatID, "Что бы посмотреть свои желания нажмите: Мой WishList")
				bot.Send(msg)
				keyboards.Menu(chatID, bot)
				return
			}

			session, err := database.CreateNewWishSession(chatID, friendID, db)
			if err != nil {
				log.Print(err)
				return
			}

			inter.FormatWishMessage(session.WishID, chatID, friendID, false, db, bot)
			return

		case text == "/menu":
			keyboard.Menu(chatID, bot)
			return

		case text == "➕ Добавить новое желание":

			inter.HandleAddNewWish(chatID, bot, db)
			return

		case text == "❌ Удалить желание":

			session, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				log.Print(err)
				return
			}

			if session.ChatID != session.TargetID {
				bot.Send(tgbot.NewMessage(chatID, "Нельзя удалять чужие желания!Ты злой и хитрый гринч)))"))
				inter.FormatWishMessage(session.WishID, chatID, session.TargetID, false, db, bot)
				return
			}

			wish, err := database.GetWishByID(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка удаления!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session.WishID, chatID, chatID, true, db, bot)
				return
			}

			err = wish.DeleteFromDB(db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка удаления!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session.WishID, chatID, chatID, true, db, bot)
				return
			}

			session.Reset()
			err = session.Save(db)
			if err != nil {
				log.Print(err)
			}

			msg := tgbot.NewMessage(chatID, "✅ *Желание Удалено!\nНадеюсь тебя порадовали новым подарком... *")
			msg.ParseMode = "Markdown"
			bot.Send(msg)

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

			session, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				log.Print(err)
				return
			}

			if session.ChatID == session.TargetID {
				bot.Send(tgbot.NewMessage(chatID, "Нельзя резервировать свои желания!Дай возможность твоим друзьям порадовать тебя!"))
				inter.FormatWishMessage(session.WishID, chatID, chatID, true, db, bot)
				return
			}

			wish, err := database.GetWishByID(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка резервации!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session.WishID, chatID, session.TargetID, false, db, bot)
				return
			}

			if wish.IsReserved {
				bot.Send(tgbot.NewMessage(chatID, "Желание уже зарезервированно другим пользователем!."))
				inter.FormatWishMessage(session.WishID, chatID, session.TargetID, false, db, bot)
				return
			}

			err = database.ReserveWish(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка резервации!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session.WishID, chatID, session.TargetID, false, db, bot)
				return
			}

			session.Reset()
			err = session.Save(db)
			if err != nil {
				log.Print(err)
			}

			msg := tgbot.NewMessage(chatID, "✅ *Желание зарезервированно!\nОбрадуйте счастливика как можно скорее! *")
			msg.ParseMode = "Markdown"
			bot.Send(msg)

		default:
			msg := tgbot.NewMessage(chatID, "Такой команды не существует!")
			bot.Send(msg)

		}

		keyboard.Menu(chatID, bot)

	}
}
