package main

import (
	"database/sql"
	"github.com/aruba8/simple-chat-logger/repository"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
)

var repo repository.Repository

func main() {
	repo = repository.InitDb()
	webhook := &tb.Webhook{
		Listen: ":" + os.Getenv("PORT"),
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: os.Getenv("WEBHOOK_URL") + "/bot" + os.Getenv("BOT_TOKEN"),
		},
	}

	settings := tb.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: webhook,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		log.Fatalf("error initializing bot: %s", err)
		return
	}

	bot.Handle(tb.OnUserJoined, func(msg *tb.Message) {

		joined := msg.UserJoined
		log.Printf("user joined. Is bot: %v, fName: %s, lName: %s, userName: %s, langCode:%s", joined.IsBot, joined.FirstName, joined.LastName, joined.Username, joined.LanguageCode)
		user := repository.User{
			TgId:         int64(joined.ID),
			FirstName:    joined.FirstName,
			LastName:     joined.LastName,
			Username:     joined.Username,
			LanguageCode: joined.LanguageCode,
			IsBot:        joined.IsBot,
		}
		_, err := repo.FindUserByTgId(user.TgId)
		if err != nil {
			if err == sql.ErrNoRows {
				id, createError := repo.CreateUser(user)
				if createError != nil {
					log.Printf("error: %v", createError)
					return
				}
				log.Printf("new user detected: %d", id)
			}
		}

	})
	//
	//bot.Handle(tb.OnText, func(msg *tb.Message) {
	//	sender := msg.Sender
	//	var user repository.User
	//	repository.DB.FirstOrCreate(&user, &repository.User{
	//		TgId:         int64(sender.ID),
	//		FirstName:    sender.FirstName,
	//		LastName:     sender.LastName,
	//		Username:     sender.Username,
	//		LanguageCode: sender.LanguageCode,
	//		IsBot:        sender.IsBot,
	//	})
	//
	//	var chat repository.Chat
	//	repository.DB.FirstOrCreate(&chat, &repository.Chat{
	//		ChatId: msg.Chat.ID,
	//		Title:  msg.Chat.Title,
	//	})
	//
	//	repository.DB.Create(&repository.Message{
	//		Datetime: time.Unix(msg.Unixtime, 0),
	//		Text:     msg.Text,
	//		ChatID:   chat.ID,
	//		UserID:   user.ID,
	//	})
	//	repository.DB.First(&user, "tg_id = ?", sender.ID)
	//	log.Printf("%v", user)
	//
	//})

	bot.Start()
}
