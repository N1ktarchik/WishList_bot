package interaction

import (
	keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ButtonProcessing(update tgbot.Update, bot *tgbot.BotAPI, msg tgbot.CallbackQuery) {

	callbackClose := tgbot.NewCallback(msg.ID, "")

	data := msg.Data
	messageID := msg.Message.MessageID
	chatID := msg.Message.Chat.ID

	defer bot.Request(callbackClose) //закрыли колл-бэк

	switch data {
	case "wishList":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		//прислать желлание из БД
		keyboard.SentKeyboard(bot, true, msg.From.ID)
		return

	case "friendsWish":
		deleteMsg := tgbot.NewDeleteMessage(chatID, messageID)
		bot.Send(deleteMsg)
		msg := tgbot.NewMessage(chatID, "To view a friend's wish list, enter the command: /friend friend_tag")
		bot.Send(msg)

	}

}
