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
	port := os.Getenv("PORT")
	webhookUrl := os.Getenv("WEBHOOK_URL")
	verbose := false
	if os.Getenv("IS_VERBOSE") == "true" {
		verbose = true
	}

	log.Print("Connecting ...\n")
	log.Printf("env:port %s\n", port)
	log.Printf("webhookUrl: %s\n", webhookUrl)

	webhook := &tb.Webhook{
		Listen: ":" + port,
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: webhookUrl,
		},
	}

	settings := tb.Settings{
		Token:   os.Getenv("BOT_TOKEN"),
		Poller:  webhook,
		Verbose: verbose,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		log.Fatalf("error initializing bot: %s", err)
		return
	}

	log.Printf("BOT initialized with url %s\n", bot.URL)
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

	bot.Handle(tb.OnText, func(msg *tb.Message) {
		log.Printf("User text: %s", msg.Text)
	})

	bot.Start()
}
