package bot

import (
	"crypto-helper/internal/infra/cache"
	"crypto-helper/internal/infra/external/coinloreApi"
	"crypto-helper/internal/services"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Handler struct {
	stateMachine *StateMachine
	service      *services.Service
}

func NewHandler() *Handler {
	botService := services.Service{
		CoinsCache: cache.NewCoinsCache(),
		UsersCache: cache.NewUsersCache(),
		Client:     &coinloreApi.Client{},
	}

	return &Handler{stateMachine: NewStateMachine(), service: &botService}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) tgbotapi.MessageConfig {

	if update.Message != nil {
		messageConfig := h.handleMessage(update.Message)
		return messageConfig
	}

	if update.CallbackQuery != nil {
		messageConfig := h.handleCallbackQuery(update.CallbackQuery)
		return messageConfig
	}

	return tgbotapi.MessageConfig{}
}

func (h *Handler) handleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	state := h.stateMachine.SetState(message.Text)

	var reply string

	switch state {

	case StateMain:
		reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
		messageConfig := h.setInlineKeyboard(message.Chat.ID, reply, []string{"перейти к списку монет"})
		return messageConfig

	case StateCoinsList:

		reply = "список всех монет:"

		coinSymbols, symbolsGettingError := h.service.GetCoinsSymbols()
		if symbolsGettingError != nil {
			log.Println(symbolsGettingError)
		}

		messageConfig := h.setCoinsKeyboard(message.Chat.ID, reply, coinSymbols)
		return messageConfig

	case StateCoinInfo:

		coinSymbol := message.Text
		coinInfo, gettingSymbolError := h.service.GetCoinInfo(coinSymbol)
		if gettingSymbolError != nil {
			log.Println(gettingSymbolError)
			return tgbotapi.NewMessage(message.Chat.ID, "ошибка получения данных")
		}

		reply := fmt.Sprintf("информация о криптовалюте %s:\nназвание: %s\nцена: %s$\nизменение цены за 1 час: %s\nизменение цены за 24 часа: %s\nизменение цены за 7 дней: %s\n", coinSymbol, coinInfo.Name, coinInfo.PriceUsd, coinInfo.PercentChange1H, coinInfo.PercentChange24H, coinInfo.PercentChange7D)

		messageConfig := h.setInlineKeyboard(message.Chat.ID, reply, []string{"назад"})

		return messageConfig

	default:
		return tgbotapi.NewMessage(message.Chat.ID, "")
	}

}

func (h *Handler) handleCallbackQuery(callback *tgbotapi.CallbackQuery) tgbotapi.MessageConfig {
	state := h.stateMachine.SetState(callback.Data)

	var reply string

	switch state {

	case StateMain:
		reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
		messageConfig := h.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"перейти к списку монет"})
		return messageConfig

	case StateCoinsList:

		reply = "список всех монет:"

		coinSymbols, symbolsGettingError := h.service.GetCoinsSymbols()
		if symbolsGettingError != nil {
			log.Println(symbolsGettingError)
		}

		messageConfig := h.setCoinsKeyboard(callback.Message.Chat.ID, reply, coinSymbols)
		return messageConfig

	case StateCoinInfo:

		coinSymbol := callback.Data
		coinInfo, gettingSymbolError := h.service.GetCoinInfo(coinSymbol)
		if gettingSymbolError != nil {
			log.Println(gettingSymbolError)
			return tgbotapi.NewMessage(callback.Message.Chat.ID, "ошибка получения данных")
		}

		reply := fmt.Sprintf("информация о криптовалюте %s:\nназвание: %s\nцена: %s$\nизменение цены за 1 час: %s\nизменение цены за 24 часа: %s\nизменение цены за 7 дней: %s\n", coinSymbol, coinInfo.Name, coinInfo.PriceUsd, coinInfo.PercentChange1H, coinInfo.PercentChange24H, coinInfo.PercentChange7D)

		messageConfig := h.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"назад"})

		return messageConfig

	default:
		return tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	}
}

func (h *Handler) setInlineKeyboard(chatId int64, text string, buttons []string) tgbotapi.MessageConfig {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	for _, button := range buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(button, button)))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	message := tgbotapi.NewMessage(chatId, text)
	message.ReplyMarkup = inlineKeyboard
	return message
}

func (h *Handler) setCoinsKeyboard(chatId int64, text string, buttons []string) tgbotapi.MessageConfig {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(buttons); i += 4 {
		end := i + 3
		if end > len(buttons) {
			end = len(buttons)
		}

		row := make([]tgbotapi.InlineKeyboardButton, end-i)
		for j := i; j < end; j++ {
			row[j-i] = tgbotapi.NewInlineKeyboardButtonData(buttons[j], buttons[j])
		}

		inlineButtons = append(inlineButtons, row)
	}

	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "назад")))

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	message := tgbotapi.NewMessage(chatId, text)
	message.ReplyMarkup = inlineKeyboard
	return message
}
