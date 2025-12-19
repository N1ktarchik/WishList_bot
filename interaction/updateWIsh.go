package interaction

import (
	"database/sql"
	"log"
	"time"

	"github.com/N1ktarchik/Wishlist_bot/database"
	"github.com/N1ktarchik/Wishlist_bot/keyboards"

	//keyboard "github.com/N1ktarchik/Wishlist_bot/keyboards"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// func HandleChangeWish(chatID int64, bot *tgbot.BotAPI, db *sql.DB) {
// 	status, err := database.GetUserStatusByID(db, chatID)
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}

// 	if status != nil && status.Step != 0 && status.IsAlive() {
// 		msg := tgbot.NewMessage(chatID,
// 			"У вас уже есть активное изменение желания.\n"+
// 				"Продолжайте вводить данные или напишите /cancel для отмены.")
// 		bot.Send(msg)
// 		return
// 	}

// 	if status != nil && !status.IsAlive() {
// 		status.Reset()
// 		err := status.Save(db)
// 		if err != nil {
// 			log.Printf("update wish done. error reset status. %v", err)
// 			msg := tgbot.NewMessage(chatID,
// 				"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
// 			bot.Send(msg)
// 			return
// 		}

// 	}

// 	status = database.CreateNewUserStatus(chatID, false, false)
// 	err = status.Save(db)
// 	if err != nil {
// 		log.Printf("update wish done. error reset status. %v", err)
// 		msg := tgbot.NewMessage(chatID,
// 			"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
// 		bot.Send(msg)
// 		return
// 	}

// 	msg := tgbot.NewMessage(chatID,
// 		"🎁 *Изменение  желания*\n\n"+
// 			"Введите *название* желания:\n"+
// 			"(или напишите /cancel для отмены)")
// 	msg.ParseMode = "Markdown"
// 	msg.ReplyMarkup = keyboard.SendNewWishAddKeyboard(bot, false, chatID)
// 	bot.Send(msg)
// }

func UpdateWishFromStatus(status *database.UserStatus, chatID int64, bot *tgbot.BotAPI, db *sql.DB, wishID int64) {
	wish := database.Wish{
		ID:          wishID,
		ChatIdLink:  status.ChatID,
		WishName:    status.WishName,
		Description: status.Description,
		Url:         status.Url,
		Price:       status.Price,
		CreatedAt:   time.Now(),
	}

	err := wish.UpdateWish(db)

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
		log.Printf("update  wish done. error reset status. %v", err)
		msg := tgbot.NewMessage(chatID,
			"save status reset error.send screenshot to admin.the operation with bot is no longer possible")
		bot.Send(msg)
		return
	}

	bot.Send(tgbot.NewMessage(chatID, "✅ Желание обновленно!"))
	keyboards.Menu(chatID, bot)
}

func CopyWishToStatus(chatID, wishID int64, db *sql.DB) (*database.UserStatus, error) {
	wish, err := database.GetWishByID(wishID, db)

	if err != nil {
		return nil, err
	}

	status := database.CreateNewUserStatus(chatID, false, false)

	status.WishName = wish.WishName
	status.Description = wish.Description
	status.Url = wish.Url
	status.Price = wish.Price
	status.Step = 5
	status.NewWish = false

	return status, nil
}
