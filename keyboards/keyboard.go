package keyboards

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Menu(update tgbot.Update, bot *tgbot.BotAPI) {
	//callback
	var (
		Keyboard = tgbot.NewInlineKeyboardMarkup(
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonData("My WishList", "wishList"),
				tgbot.NewInlineKeyboardButtonData("Friends WishList", "friendsWish"),
			),
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonURL("Info", "https://music.yandex.ru"), //добавить статью телеграф (мануал по пользованию)
				tgbot.NewInlineKeyboardButtonURL("Help", "https://t.me/n1ktarchik"),
			),
			tgbot.NewInlineKeyboardRow(
				//tgbot.NewInlineKeyboardButtonURL("The news channel", "https://music.yandex.ru"),   //канал тгк (добавить ссылку)
				tgbot.NewInlineKeyboardButtonURL("Support the author", "https://music.yandex.ru"), //запрос деняк (добавить ссылку)
			),
		)
	)

	msg := tgbot.NewMessage(update.Message.Chat.ID, "Choose an action")
	msg.ReplyMarkup = Keyboard
	bot.Send(msg)

}

func SentKeyboard(bot *tgbot.BotAPI, choise bool, chatid int64) {
	//text
	//choise=true
	keyboardToMyWish := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(

			tgbot.NewKeyboardButton("➕ Add new wish"),
			tgbot.NewKeyboardButton("❌ Delete wish"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("✏️ Change wish"),
			tgbot.NewKeyboardButton("➡️ Next wish"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("🔙 Exit to main menu"),
		),
	)

	keyboardToMyWish.ResizeKeyboard = true
	keyboardToMyWish.OneTimeKeyboard = true
	keyboardToMyWish.Selective = true

	//choise=false
	keyboardToFriendWish := tgbot.NewReplyKeyboard(

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("✅ Reserve wish"),
			tgbot.NewKeyboardButton("➡️ Next wish"),
		),

		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton("🔙 Exit to main menu"),
		),
	)

	keyboardToFriendWish.ResizeKeyboard = true
	keyboardToFriendWish.OneTimeKeyboard = true
	keyboardToFriendWish.Selective = true

	msg := tgbot.NewMessage(chatid, "ㅤ") //корейский "чистый" пробел

	if choise {
		msg.ReplyMarkup = keyboardToMyWish
	} else {
		msg.ReplyMarkup = keyboardToFriendWish
	}

	bot.Send(msg)
}
