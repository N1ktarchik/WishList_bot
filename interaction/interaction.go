package interaction

import (
	"database/sql"
	"fmt"
	"log"

	"strings"

	"github.com/N1ktarchik/Wishlist_bot/database"
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
			return
		}

		FormatWishMessage(session.WishID, chatID, chatID, true, db, bot)

	case "friendsWish":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		msg := tgbot.NewMessage(chatID, "Чтобы просмотреть список желаний друга, введите команду: /friend friend_tag")
		bot.Send(msg)

	case "donate":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		sendDonateMessage(bot, chatID)

	}

}

func FormatWishMessage(wishID, chatID, wishOwner int64, isOwnWish bool, db *sql.DB, bot *tgbot.BotAPI) {

	var builder strings.Builder
	wish, err := database.GetWishByID(wishID, db)
	if err != nil {
		log.Print(err)
		bot.Send(tgbot.NewMessage(wishOwner, "У пользователя нет желаний!"))
	}

	if isOwnWish {
		builder.WriteString("📋 *Мое желание*\n\n")
	} else {
		username, err := database.GetUsernameByID(wishOwner, db)
		if err != nil {
			log.Print(err)
			bot.Send(tgbot.NewMessage(chatID, "Что-то пошло не так...\nОтправь скриншот в поддержку или попробуй снова. Error to get wish by id"))
			return
		}
		builder.WriteString(fmt.Sprintf("🎁 *Желание @%s*\n\n", username))
	}

	builder.WriteString(fmt.Sprintf("📌 *Название:* %s\n", wish.WishName))

	if wish.Description != "" {
		builder.WriteString(fmt.Sprintf("📝 *Описание:* %s\n", wish.Description))
	}

	if wish.Price > 0 {
		builder.WriteString(fmt.Sprintf("💰 *Цена:* %.2f руб.\n", wish.Price))
	}

	if wish.Url != "" {
		builder.WriteString(fmt.Sprintf("🔗 *Ссылка:* %s\n", wish.Url))
	}

	flag := false
	if wish.IsReserved {
		if !isOwnWish {
			builder.WriteString("\n")
			builder.WriteString("🚫 *ЗАРЕЗЕРВИРОВАНО!*\n")
			builder.WriteString("_(Это желание уже забронировано другим пользователем)_\n")
			flag = true
		}
	} else if !isOwnWish {
		builder.WriteString("\n")
		builder.WriteString("✅ *Доступно для резервирования*\n")
	}

	message := builder.String()

	msg := tgbot.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	bot.Send(msg)

	if flag {
		keyboard.SentWishReservedKeyboard(bot, chatID)
		return
	}

	keyboard.SentWishKeyboard(bot, isOwnWish, chatID)

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

	if _, err := bot.Send(msg); err != nil {
		log.Printf("error send donate message: %v", err)
	}
}
