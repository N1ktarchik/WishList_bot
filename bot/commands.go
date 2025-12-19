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

	//for beta test
	if !database.ChekTesterRights(update.Message.Chat.ID, db) {

		if update.Message != nil && update.Message.Text != "" && strings.HasPrefix(update.Message.Text, "/test") {
			if inter.CheckPassword(update.Message.Text) {
				database.SaveNewTester(update.Message.Chat.ID, db)
				bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Доступ к тестировке разрешен!\nСпасибо за твою помощь в моем проекте!"))
				keyboard.SendTesterKeyboard(bot, update.Message.Chat.ID)
				return
			}

		}
		bot.Send(tgbot.NewMessage(update.Message.Chat.ID,
			"Воспользоваться ботом сейчас не получится.\n\nСейчас проходят тесты бота и вносятся правки.\n\nЕсли хочешь получить доступ к бета-тестированию, заходи в тгк, и читай последний пост.\n\nhttps://t.me/n1k_go"))
		bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Если ты уже в команде тестировщиков, введи команду:\n\n/test password\n\nВместо password укажи пароль тестировщика"))
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
			bot.Send(tgbot.NewMessage(chatID, "⏳ Время добавления желания истекло. Начните заново."))
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
				msg := tgbot.NewMessage(chatID,
					"Тег не распознается! Пожалуйста, отправьте его таким образом:\n\n"+
						"`/friend @username`")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			}

			friendID, err := database.GetIdByUsername(mas[1], db)
			if err != nil {
				msg := tgbot.NewMessage(chatID, fmt.Sprintf("Пользователь с username %s не найден!", mas[1]))
				log.Print(err)
				bot.Send(msg)
				keyboards.Menu(chatID, bot)
				return
			}

			if friendID == chatID {
				msg := tgbot.NewMessage(chatID, "Что бы посмотреть свои желания нажмите: Мой WishList")
				bot.Send(msg)
				keyboards.Menu(chatID, bot)
				return
			}

			session, err := database.CreateNewWishSession(chatID, friendID, db)
			if err != nil {
				log.Print(err)
				keyboards.Menu(chatID, bot)
				return
			}

			inter.FormatWishMessage(session, db, bot)
			return

		case text == "➕ Добавить новое желание":

			inter.HandleAddNewWish(chatID, bot, db)
			return

		case text == "❌ Удалить желание":

			session, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				log.Print(err)
				keyboards.Menu(chatID, bot)
				return
			}

			if session.ChatID != session.TargetID {
				bot.Send(tgbot.NewMessage(chatID, "Нельзя удалять чужие желания!Ты злой и хитрый гринч)))"))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			wish, err := database.GetWishByID(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка удаления!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			err = wish.DeleteFromDB(db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка удаления!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session, db, bot)
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
			session, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				log.Print(err)
				keyboards.Menu(chatID, bot)
				return
			}

			if session.ChatID != session.TargetID {
				bot.Send(tgbot.NewMessage(chatID, "Нельзя изменять чужие желания!Ты злой и хитрый гринч)))"))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			session.UpdateLiveTime(30)

			inter.HandleChangeWish(chatID, bot, db)
			return

		case text == "➡️ Следующее желание":
			inter.ScrollingWish(bot, chatID, true, db)
			return

		case text == "⬅️ Предыдущее желание":
			inter.ScrollingWish(bot, chatID, false, db)
			return

		case text == "🔙 Вернуться в главное меню" || text == "/menu":
			keyboard.Menu(chatID, bot)
			return

		case text == "✅ Зарезервировать желание":

			session, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				log.Print(err)
				keyboards.Menu(chatID, bot)
				return
			}

			if session.ChatID == session.TargetID {
				bot.Send(tgbot.NewMessage(chatID, "Нельзя резервировать свои желания!Дай возможность твоим друзьям порадовать тебя!"))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			wish, err := database.GetWishByID(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка резервации!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			if wish.IsReserved {
				bot.Send(tgbot.NewMessage(chatID, "Желание уже зарезервированно другим пользователем!."))
				inter.FormatWishMessage(session, db, bot)
				return
			}

			err = database.ReserveWish(session.WishID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка резервации!Попробуйте снова или обратитесь в поддержку."))
				inter.FormatWishMessage(session, db, bot)
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
