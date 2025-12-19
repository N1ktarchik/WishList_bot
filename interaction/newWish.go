package interaction

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/N1ktarchik/Wishlist_bot/database"
	"github.com/N1ktarchik/Wishlist_bot/keyboards"
	keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleAddNewWish(chatID int64, bot *tgbot.BotAPI, db *sql.DB) {
	status, err := database.GetUserStatusByID(db, chatID)
	if err != nil {
		log.Print(err)
		return
	}

	if status != nil && status.Step != 0 && status.IsAlive() {
		msg := tgbot.NewMessage(chatID,
			"У вас уже есть активное добавление.\n"+
				"Продолжайте вводить данные или напишите /cancel для отмены.")
		bot.Send(msg)
		return
	}

	if status != nil && !status.IsAlive() {
		status.Reset()
		err := status.Save(db)
		if err != nil {
			log.Printf("save new wish done. error reset status. %v", err)
			msg := tgbot.NewMessage(chatID,
				"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
			bot.Send(msg)
			return
		}

	}

	status = database.CreateNewUserStatus(chatID, true)
	err = status.Save(db)
	if err != nil {
		log.Printf("save new wish done. error reset status. %v", err)
		msg := tgbot.NewMessage(chatID,
			"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
		bot.Send(msg)
		return
	}

	msg := tgbot.NewMessage(chatID,
		"🎁 *Добавление нового желания*\n\n"+
			"Введите *название* желания:\n"+
			"(или напишите /cancel для отмены)")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard.SendNewWishAddKeyboard(bot, false, chatID)
	bot.Send(msg)
}

func ProcessingNewWish(status *database.UserStatus, update tgbot.Update, bot *tgbot.BotAPI, db *sql.DB) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	txt := strings.TrimSpace(update.Message.Text)

	if txt == "/cancel" || txt == "❌ Отмена" {
		status.Delete(db)
		bot.Send(tgbot.NewMessage(chatID, "✅ Добавление отменено"))
		keyboard.Menu(chatID, bot)
		return
	}

	status.UpdateLiveTime(5)
	switch status.Step {

	case 1:
		if txt == "" {
			msg := tgbot.NewMessage(chatID, "Название желания не может быть пустым! Попробуйте еще раз:")
			msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, false, chatID)
			bot.Send(msg)
			return
		}

		if len(txt) < 3 {
			msg := tgbot.NewMessage(chatID, "Название желания слишком короткое!  Попробуйте еще раз:")
			msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, false, chatID)
			bot.Send(msg)
			return
		}

		status.WishName = txt
		status.Step = 2

		err := status.Save(db)
		if err != nil {
			bot.Send(tgbot.NewMessage(chatID,
				"Ошибка сохранения! Попробуйте еще раз, или пришлите скриншот в поддержку!"))
			return
		}

		msg := tgbot.NewMessage(chatID,
			"Введите *описание* желания:")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
		bot.Send(msg)

	case 2:
		if txt == "🚫 Пропустить" {
			txt = ""
		}

		if len(txt) > 1000 {
			msg := tgbot.NewMessage(chatID,
				"❌ Описание слишком длинное! Максимум 1000 символов.")
			msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
			bot.Send(msg)
			return
		}

		status.Step = 3
		status.Description = txt

		err := status.Save(db)
		if err != nil {
			bot.Send(tgbot.NewMessage(chatID,
				"Ошибка сохранения! Попробуйте еще раз, или пришлите скриншот в поддержку!"))
			return
		}

		msg := tgbot.NewMessage(chatID,
			"Введите *ссылку на товар*:")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
		bot.Send(msg)

	case 3:
		if txt != "🚫 Пропустить" {
			if !strings.HasPrefix(txt, "http://") && !strings.HasPrefix(txt, "https://") {
				msg := tgbot.NewMessage(chatID,
					"❌ Ссылка должна начинаться с http:// или https://\n"+
						"Попробуйте еще раз:")
				msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
				bot.Send(msg)
				return
			}

			if len(txt) > 1000 {
				msg := tgbot.NewMessage(chatID,
					"❌ Ссылка слишком длинная! Максимум 1000 символов.")
				msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
				bot.Send(msg)
				return
			}
		} else {
			txt = ""
		}

		status.Step = 4
		status.Url = txt

		err := status.Save(db)
		if err != nil {
			bot.Send(tgbot.NewMessage(chatID,
				"Ошибка сохранения! Попробуйте еще раз, или пришлите скриншот в поддержку!"))
			return
		}

		msg := tgbot.NewMessage(chatID,
			"Введите *цену* в рублях (только число, например: 1500.50):")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
		bot.Send(msg)

	case 4:
		var price float64 = 0

		if txt != "🚫 Пропустить" {

			if strings.Contains(txt, ",") {
				txt = strings.ReplaceAll(txt, ",", ".")
			}

			if !strings.Contains(txt, ".") {
				txt = txt + ".00"
			}

			parsedPrice, err := strconv.ParseFloat(txt, 64)

			if err != nil {
				msg := tgbot.NewMessage(chatID, "❌ Цена некорректна, попробуйте еще раз:")
				msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
				bot.Send(msg)
				return
			}

			if parsedPrice <= 0 {
				msg := tgbot.NewMessage(chatID,
					"❌ Цена не может быть меньше или равно нулю, попробуйте еще раз:")
				msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, true, chatID)
				bot.Send(msg)
				return
			}

			price = parsedPrice

		}

		status.Price = price
		status.Step = 5

		err := status.Save(db)
		if err != nil {
			bot.Send(tgbot.NewMessage(chatID,
				"Ошибка сохранения! Попробуйте еще раз, или пришлите скриншот в поддержку!"))
			return
		}

		SendConfirmation(status, chatID, bot)
		status.Step = 5

	case 5:
		HandleConfirmation(status, txt, chatID, bot, db)
		return

	default:
		log.Printf("Неизвестный шаг статуса: %d для пользователя %d", status.Step, userID)
		status.Reset()
		status.Save(db)

		msg := tgbot.NewMessage(chatID, "⚠️ Произошла ошибка. Начинаем заново.")
		msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, false, chatID)
		bot.Send(msg)

	}
}

func SendConfirmation(status *database.UserStatus, chatID int64, bot *tgbot.BotAPI) {
	msgText := fmt.Sprintf(
		"🎯 *Проверьте данные:*\n\n"+
			"📝 *Название:* %s\n",
		status.WishName)

	if status.Price != 0 {
		msgText += fmt.Sprintf("💰 *Цена:* %.2f руб.\n", status.Price)
	} else {
		msgText += "💰 *Цена:* не указана\n"
	}

	if status.Description != "" {
		msgText += fmt.Sprintf("📋 *Описание:* %s\n", status.Description)
	} else {
		msgText += "📋 *Описание:* не указано\n"
	}

	if status.Url != "" {
		msgText += fmt.Sprintf("🔗 *Ссылка:* %s\n", status.Url)
	} else {
		msgText += "🔗 *Ссылка:* не указана\n"
	}

	msgText += "\n*Всё верно?*"

	msg := tgbot.NewMessage(chatID, msgText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard.SendConfirmationKeyboard(bot, chatID)
	bot.Send(msg)
}

func HandleConfirmation(status *database.UserStatus, text string, chatID int64, bot *tgbot.BotAPI, db *sql.DB) {
	switch strings.ToLower(text) {
	case "да", "yes", "ok", "подтверждаю", "✅ да! сохранить.":
		if status.NewWish {
			SaveWishFromStatus(status, chatID, bot, db)
		} else {
			wish, err := database.GetWishSessonByID(chatID, db)
			if err != nil {
				bot.Send(tgbot.NewMessage(chatID, "Ошибка обновления желания!Попробуйте сновы или обратитесь в поддержку"))
				status.Reset()
				status.Save(db)
				log.Print(err)
				return
			}

			UpdateWishFromStatus(status, chatID, bot, db, wish.WishID)
		}

	case "нет", "no", "отмена", "❌ нет! начать заново.":
		status.Step = 1
		status.Save(db)
		msg := tgbot.NewMessage(chatID, "Введите название:")
		msg.ReplyMarkup = keyboards.SendNewWishAddKeyboard(bot, false, chatID)
		bot.Send(msg)

	default:
		msg := tgbot.NewMessage(chatID, "Пожалуйста, выберите вариант из клавиатуры:")
		msg.ReplyMarkup = keyboards.SendConfirmationKeyboard(bot, chatID)
		bot.Send(msg)
	}
}

func SaveWishFromStatus(status *database.UserStatus, chatID int64, bot *tgbot.BotAPI, db *sql.DB) {
	wish := database.Wish{
		ChatIdLink:  status.ChatID,
		WishName:    status.WishName,
		Description: status.Description,
		Url:         status.Url,
		Price:       status.Price,
		IsReserved:  false,
		CreatedAt:   time.Now(),
	}

	err := wish.AddToDB(db)

	if err != nil {
		log.Printf("error save wish: %v", err)
		msg := tgbot.NewMessage(chatID,
			"❌ Ошибка сохранения. Попробуйте еще раз или пришлите скриншот в поддержку.")
		bot.Send(msg)
		status.Reset()
		err := status.Save(db)
		if err != nil {
			log.Printf("error save new wish. error reset status. %v", err)
		}
		return
	}

	status.Reset()
	err = status.Save(db)

	if err != nil {
		log.Printf("save new wish done. error reset status. %v", err)
		msg := tgbot.NewMessage(chatID,
			"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
		bot.Send(msg)
		return
	}

	bot.Send(tgbot.NewMessage(chatID, "✅ Желание сохранено!"))
	keyboards.Menu(chatID, bot)
}
