package main

import (
	"crypto-helper/internal/handler"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	envLoadingError := godotenv.Load()
	if envLoadingError != nil {
		log.Fatalf("Error loading .env file: %s", envLoadingError)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	telegramBot, creatingBotError := handler.NewBot(token)
	if creatingBotError != nil {
		log.Fatalf("Error creating bot: %s", creatingBotError)
	}

	telegramBot.Start()
}
