package main

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
)

func main() {
	// Use config.yml to store the bot token and server URL
	config := getConfig()
	botToken := config.Telegram.BotToken
	serverUrl := config.Telegram.ServerUrl
	webhookPath := config.Telegram.WebhookPath

	bot, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		log.Println("Error creating Telegram bot:", err)
		return
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbot.NewWebhook(serverUrl + webhookPath))
	if err != nil {
		log.Println("Error setting webhook:", err)
		return
	}

	updates := bot.ListenForWebhook(webhookPath)

	// Start HTTP server to listen for webhooks
	err = http.ListenAndServe("0.0.0.0:8081", http.DefaultServeMux)
	if err != nil {
		log.Println("Error starting HTTP server:", err)
		return
	}

	// Loop to handle incoming messages
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			botResponse := "hello!"
			msg := tgbot.NewMessage(update.Message.Chat.ID, botResponse)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Error sending message:", err)
			}
		}
	}
}
