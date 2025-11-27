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

// func ChoiseToMyWishListProcessing(choise string, update tgbot.Update, bot *tgbot.BotAPI) {

// 	switch choise {
// 	case "➕ Add new wish":
// 		//обработка через БД

// 	case "❌ Delete wish":
// 		//обработка через БД

// 	case "✏️ Change wish":
// 		//обработка через БД

// 	case "➡️ Next wish":
// 		//обработка через БД
// 		//Обдумать архитектуру

// 	case "🔙 Exit to main menu":
// 		keyboard.Menu(update, bot)
// 		return

// 	default:
// 		msg := tgbot.NewMessage(update.Message.Chat.ID, "The command is not recognized. Select the command on the keyboard 👇")
// 		bot.Send(msg)
// 		keyboard.SentKeyboard(bot, true, update.Message.Chat.ID)

// 	}
// }
