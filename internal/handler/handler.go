package handler

import (
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/infra/external/coinloreApi"
	"crypto-helper/internal/services"
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

	updatingError := b.handler.service.UpdateCoinsInfo()
	if updatingError != nil {
		log.Printf("Error updating coins info: %s", updatingError.Error())
	}
	b.handler.service.StartUpdatingCoinsInfo()

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
			log.Printf("[%d] %s", update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

			messageConfig := b.handler.HandleUpdate(update)
			_, sendingError := b.api.Send(messageConfig)
			if sendingError != nil {
				log.Printf("Error sending message: %s", sendingError)
			}
		}
	}
}

type Handler struct {
	stateMachine *services.StateMachine
	service      *services.Service
}

var currentCoin string

func NewHandler() *Handler {
	botService := services.Service{
		CoinsCache:   cache.NewCoinsCache(),
		UsersCache:   cache.NewUsersCache(),
		Client:       &coinloreApi.Client{},
		StateMachine: services.NewStateMachine(),
	}

	return &Handler{stateMachine: services.NewStateMachine(), service: &botService}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) tgbotapi.MessageConfig {
	if update.Message != nil {
		messageConfig := h.service.HandleMessage(*update.Message, &currentCoin)
		return messageConfig
	}

	if update.CallbackQuery != nil {
		messageConfig := h.service.HandleCallbackQuery(update.CallbackQuery, &currentCoin)
		return messageConfig
	}

	return tgbotapi.MessageConfig{}
}
