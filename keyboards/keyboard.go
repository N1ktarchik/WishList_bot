package keyboards

import (
	"log"
	"time"

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
				tgbot.NewInlineKeyboardButtonURL("Support the author", "https://music.yandex.ru"), //запрос деняк (добавить ссылку)
			),
		)
	)

	msg := tgbot.NewMessage(chatID, "Выберите действие")
	msg.ReplyMarkup = Keyboard
	bot.Send(msg)

}

func SentWishKeyboard(bot *tgbot.BotAPI, choise bool, chatid int64) {

	msg := tgbot.NewMessage(chatid, "Обрабатываю команду...")
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	time.Sleep(time.Millisecond * 300)
	deleteMsg := tgbotapi.NewDeleteMessage(chatid, sentMsg.MessageID)
	bot.Send(deleteMsg)

	//text
	//choise=true
	keyboardToMyWish := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("➕ Добавить новое желание"),
			tgbot.NewKeyboardButton("❌ Удалить желание"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("⬅️ Предыдущие желание"),
			tgbot.NewKeyboardButton("➡️ Следующее желание"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("✏️ Изменить желание"),
			tgbot.NewKeyboardButton("🔙 Вернуться в главное меню"),
		),
	)

	keyboardToMyWish.ResizeKeyboard = true
	keyboardToMyWish.OneTimeKeyboard = true
	keyboardToMyWish.Selective = true

	//choise=false
	keyboardToFriendWish := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("✅ Зарезервировать желание"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("➡️ Следующее желание"),
			tgbot.NewKeyboardButton("⬅️ Предыдущие желание"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("🔙 Вернуться в главное меню"),
		),
	)

	keyboardToFriendWish.ResizeKeyboard = true
	keyboardToFriendWish.OneTimeKeyboard = true
	keyboardToFriendWish.Selective = true

	sms := tgbot.NewMessage(chatid, "Выбери команду на предложенной клавиатуре: ")

	if choise {
		sms.ReplyMarkup = keyboardToMyWish
	} else {
		sms.ReplyMarkup = keyboardToFriendWish
	}

	bot.Send(sms)
}

func SentNewWishAddKeyboard(bot *tgbot.BotAPI, choise bool, chatid int64) *tgbot.ReplyKeyboardMarkup {

	//true=with Skip
	KeyboardWithSkip := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("❌ Отмена"),
			tgbot.NewKeyboardButton("🚫 Пропустить"),
		),
	)

	KeyboardWithSkip.ResizeKeyboard = true
	KeyboardWithSkip.OneTimeKeyboard = true
	KeyboardWithSkip.Selective = true

	if choise {
		return &KeyboardWithSkip
	}

	//false = with out Skip
	KeyboardWithOutSkip := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("❌ Отмена"),
		),
	)

	KeyboardWithOutSkip.ResizeKeyboard = true
	KeyboardWithOutSkip.OneTimeKeyboard = true
	KeyboardWithOutSkip.Selective = true

	return &KeyboardWithOutSkip
}

func SentConfirmationKeyboard(bot *tgbot.BotAPI, chatid int64) *tgbot.ReplyKeyboardMarkup {
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
