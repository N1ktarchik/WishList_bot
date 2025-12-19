package interaction

import (
	"database/sql"
	"fmt"
	"log"

	"strings"

	"github.com/N1ktarchik/Wishlist_bot/database"
	"github.com/N1ktarchik/Wishlist_bot/keyboards"
	keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ButtonProcessing(update tgbot.Update, bot *tgbot.BotAPI, msg tgbot.CallbackQuery, db *sql.DB) {

	callbackClose := tgbot.NewCallback(msg.ID, "")

	data := msg.Data
	messageID := msg.Message.MessageID
	chatID := msg.Message.Chat.ID

	defer bot.Request(callbackClose) //закрыли колл-бэк

	switch data {
	case "wishList":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)

		session, err := database.CreateNewWishSession(chatID, chatID, db)
		if err != nil {
			bot.Send(tgbot.NewMessage(chatID, "Твой список желаний пуст!"))
			keyboards.SendFirstWishKeyboard(bot, chatID)

			return
		}

		FormatWishMessage(session, db, bot)

	case "friendsWish":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		msg := tgbot.NewMessage(chatID,
			"Чтобы просмотреть список желаний друга, введите команду:\n\n"+
				"`/friend @username`")
		msg.ParseMode = "Markdown"
		bot.Send(msg)

	case "donate":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		sendDonateMessage(bot, chatID)

	}

}

func FormatWishMessage(session *database.WishSession, db *sql.DB, bot *tgbot.BotAPI) {

	var builder strings.Builder
	wish, err := database.GetWishByID(session.WishID, db)
	if err != nil {
		log.Print(err)
		bot.Send(tgbot.NewMessage(session.TargetID, "У пользователя нет желаний!"))
	}

	ownWish := session.ChatID == session.TargetID
	if ownWish {
		builder.WriteString("📋 *Мое желание*\n\n")
	} else {
		username, err := database.GetUsernameByID(session.TargetID, db)
		if err != nil {
			log.Print(err)
			bot.Send(tgbot.NewMessage(session.ChatID, "Что-то пошло не так...\nОтправь скриншот в поддержку или попробуй снова. Error to get wish by id"))
			return
		}
		builder.WriteString(fmt.Sprintf("🎁 *Желание @%s*\n\n", username))
	}

	builder.WriteString(fmt.Sprintf("📌 *Название:* %s\n", wish.WishName))

	if wish.Description != "" {
		builder.WriteString(fmt.Sprintf("📝 *Описание:* %s\n", wish.Description))
	} else {
		builder.WriteString("📝 *Описание:* _не указано _\n")
	}

	if wish.Price > 0 {
		builder.WriteString(fmt.Sprintf("💰 *Цена:* %.2f руб.\n", wish.Price))
	} else {
		builder.WriteString("💰 *Цена:* _не указана _\n")
	}

	if wish.Url != "" {
		builder.WriteString(fmt.Sprintf("🔗 *Ссылка:* %s\n", wish.Url))
	} else {
		builder.WriteString("🔗 *Ссылка:* _не указана _\n")
	}

	reserved := false
	if wish.IsReserved {
		if !ownWish {
			builder.WriteString("\n")
			builder.WriteString("🚫 *ЗАРЕЗЕРВИРОВАНО!*\n")
			builder.WriteString("_(Это желание уже забронировано другим пользователем)_\n")
			reserved = true
		}
	} else if !ownWish {
		builder.WriteString("\n")
		builder.WriteString("✅ *Доступно для резервирования*\n")
	}

	message := builder.String()

	msg := tgbot.NewMessage(session.ChatID, message)
	msg.ParseMode = "Markdown"

	bot.Send(msg)

	navigation, _ := database.GetWishNavigation(session.WishID, session.TargetID, db)
	keyboard.SendWishKeyboard(bot, ownWish, session.ChatID, navigation, reserved)
}

func sendDonateMessage(bot *tgbotapi.BotAPI, chatID int64) {
	messageText := `🎁 <b>Поддержать проект | Финансовая помощь автору</b>

	<u>Все средства идут исключительно на:</u>
	• 🖥 Оплату хостинга и серверов
	• 🔄 Обновления и поддержку бота
	• 🛠 Разработку новых функций
	• 📊 Мониторинг и аналитику

	════════════════════════════

	<b>💳 ВАРИАНТЫ ОПЛАТЫ:</b>

	<b>1. 🏦 Для клиентов Т-банка</b> — <i>самый удобный для меня и вас</i>
	<a href="https://tbank.ru/cf/PeWKHqZMRp">https://tbank.ru/cf/PeWKHqZMRp</a>

	<b>2. 💎 Перевод на карту</b> — <i>с любого российского банка</i>
	<code>2200 7013 3782 4293</code>
	<i>Имя получателя: Никита.К </i>

	<b>3. 🔄 Другие банки</b> — <i>через сторонний сервис</i>
	⚠️  <b>Внимание:</b> высокая комиссия для меня (до 15%)
	<a href="https://pay.cloudtips.ru/p/23cddc84">https://pay.cloudtips.ru/p/23cddc84</a>

	════════════════════════════

	<b>🌟 БОЛЬШОЕ СПАСИБО ЗА ВАШУ ПОДДЕРЖКУ! 🌟</b>

	<i>Ваш вклад позволяет боту:</i>
	✓ Работать 24/7 без перерывов
	✓ Быстро отвечать на запросы
	✓ Обрабатывать больше пользователей
	✓ Развиваться и становиться лучше

	💬 <i>"Даже маленькая помощь — большой шаг вперёд"</i>

	🤖 <b>Спасибо, что вы с нами!</b>`

	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("error send donate message: %v", err)
	}

	keyboards.SendBackMainMenuKeyboard(bot, chatID)
}

func ScrollingWish(bot *tgbotapi.BotAPI, chatID int64, next bool, db *sql.DB) error {
	session, err := database.GetWishSessonByID(chatID, db)

	if err != nil {
		if err == sql.ErrNoRows {
			bot.Send(tgbot.NewMessage(chatID, "Выберите чей виш лист вы хотите посмотреть!"))
			keyboards.Menu(chatID, bot) //отправка меню для пользователя
			return nil
		}

		keyboards.Menu(chatID, bot)
		return err
	}

	if session == nil {
		bot.Send(tgbot.NewMessage(chatID, "Выберите чей виш лист вы хотите посмотреть!"))
		keyboards.Menu(chatID, bot)
		return nil
	}

	navigation, err := database.GetWishNavigation(session.WishID, session.TargetID, db)

	if err != nil {
		FormatWishMessage(session, db, bot)
		return err
	}

	change := false

	if next && navigation.NextID != nil {
		session.WishID = *navigation.NextID
		change = true
	} else if !next && navigation.PrevID != nil {
		session.WishID = *navigation.PrevID
		change = true
	}

	if !change {

		if next {
			bot.Send(tgbot.NewMessage(chatID, "Это было последнее желание пользователя!"))
		} else {
			bot.Send(tgbot.NewMessage(chatID, "Это было само первое желание пользователя!"))
		}

		FormatWishMessage(session, db, bot)
		return nil
	}

	session.UpdateLiveTime(10)
	err = session.Update(db)
	if err != nil {
		log.Printf("error to update session. %v", err)
	}
	FormatWishMessage(session, db, bot)
	return nil

}
