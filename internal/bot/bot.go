package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	handler := NewHandler()

	return &Bot{api: api, handler: handler}, nil
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	updatingErr := b.handler.service.UpdateCoinsInfo()
	if updatingErr != nil {
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			messageConfig := b.handler.HandleUpdate(update)
			_, sendingError := b.api.Send(messageConfig)
			if sendingError != nil {
				log.Printf("Error sending message: %s", sendingError)
			}
		}

		if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

			messageConfig := b.handler.HandleUpdate(update)
			_, sendingError := b.api.Send(messageConfig)
			if sendingError != nil {
				log.Printf("Error sending message: %s", sendingError)
			}
		}
	}

}
