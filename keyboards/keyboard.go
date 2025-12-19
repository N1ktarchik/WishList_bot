package keyboards

import (
	db "github.com/N1ktarchik/Wishlist_bot/database"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Menu(chatID int64, bot *tgbot.BotAPI) {
	//callback
	var (
		Keyboard = tgbot.NewInlineKeyboardMarkup(
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonData("Мой WishList", "wishList"),
				tgbot.NewInlineKeyboardButtonData("WishList друга", "friendsWish"),
			),
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonURL("Info", "https://t.me/n1k_go"), //тгк (возможно замена на статью телеграф)
				tgbot.NewInlineKeyboardButtonURL("Help", "https://t.me/n1ktarchik"),
			),
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonData("Support the author", "donate"), //запрос деняк (добавить ссылку)
			),
		)
	)

	msg := tgbot.NewMessage(chatID, "Выберите действие")
	msg.ReplyMarkup = Keyboard
	bot.Send(msg)

}

func SendWishKeyboard(bot *tgbotapi.BotAPI, isMyWish bool, chatID int64, navigation *db.WishNavigation, reserved bool) {
	var rows [][]tgbotapi.KeyboardButton

	if isMyWish {
		rows = append(rows, []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("➕ Добавить новое желание"),
			tgbotapi.NewKeyboardButton("❌ Удалить желание"),
		})

		addNavigationButtons(&rows, navigation)

		rows = append(rows, []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("✏️ Изменить желание"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		})
	} else {
		addReserveButton(&rows, reserved)

		addNavigationButtons(&rows, navigation)

		rows = append(rows, []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		})
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true
	keyboard.Selective = true

	msg := tgbotapi.NewMessage(chatID, "Выбери команду:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func addReserveButton(rows *[][]tgbotapi.KeyboardButton, reserved bool) {
	var navButtons []tgbotapi.KeyboardButton

	if !reserved {
		navButtons = append(navButtons, tgbotapi.NewKeyboardButton("✅ Зарезервировать желание"))
	}

	if len(navButtons) > 0 {
		*rows = append(*rows, navButtons)
	}
}

func addNavigationButtons(rows *[][]tgbotapi.KeyboardButton, navigation *db.WishNavigation) {

	if navigation == nil {
		return
	}

	var navButtons []tgbotapi.KeyboardButton

	if navigation.PrevID != nil {
		navButtons = append(navButtons, tgbotapi.NewKeyboardButton("⬅️ Предыдущее желание"))
	}

	if navigation.NextID != nil {
		navButtons = append(navButtons, tgbotapi.NewKeyboardButton("➡️ Следующее желание"))
	}

	if len(navButtons) > 0 {
		*rows = append(*rows, navButtons)
	}
}

func SendNewWishAddKeyboard(bot *tgbot.BotAPI, withSkip bool, chatid int64) *tgbot.ReplyKeyboardMarkup {

	var Keyboard = tgbot.NewReplyKeyboard()

	if withSkip {
		Keyboard = tgbot.NewReplyKeyboard(

			tgbot.NewKeyboardButtonRow(

				tgbot.NewKeyboardButton("❌ Отмена"),
				tgbot.NewKeyboardButton("🚫 Пропустить"),
			),
		)
	} else {
		Keyboard = tgbot.NewReplyKeyboard(

			tgbot.NewKeyboardButtonRow(

				tgbot.NewKeyboardButton("❌ Отмена"),
			),
		)
	}

	Keyboard.ResizeKeyboard = true
	Keyboard.OneTimeKeyboard = true
	Keyboard.Selective = true

	return &Keyboard
}

func SendConfirmationKeyboard(bot *tgbot.BotAPI, chatid int64) *tgbot.ReplyKeyboardMarkup {
	Keyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("✅ Да! Сохранить."),
			tgbot.NewKeyboardButton("❌ Нет! Начать заново."),
		),
	)

	Keyboard.ResizeKeyboard = true
	Keyboard.OneTimeKeyboard = true
	Keyboard.Selective = true

	return &Keyboard

}

func SendFirstWishKeyboard(bot *tgbot.BotAPI, chatid int64) {
	keyboard := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("➕ Добавить новое желание"),
		),
	)

	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true
	keyboard.Selective = true

	sms := tgbot.NewMessage(chatid, "Хоешь добавить первое желание? ")
	sms.ReplyMarkup = keyboard
	bot.Send(sms)
}

func SendBackMainMenuKeyboard(bot *tgbot.BotAPI, chatid int64) {
	keyboard := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("🔙 Вернуться в главное меню"),
		),
	)

	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true
	keyboard.Selective = true

	sms := tgbot.NewMessage(chatid, " ❤️ ")
	sms.ReplyMarkup = keyboard
	bot.Send(sms)
}
