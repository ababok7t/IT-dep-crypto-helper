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

var currentCoin string

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
		messageConfig := h.handleMessage(*update.Message, &currentCoin)
		return messageConfig
	}

	if update.CallbackQuery != nil {
		messageConfig := h.handleCallbackQuery(update.CallbackQuery, &currentCoin)
		return messageConfig
	}

	return tgbotapi.MessageConfig{}
}

func (h *Handler) handleMessage(message tgbotapi.Message, currentCoin *string) tgbotapi.MessageConfig {
	state := h.stateMachine.SetState(message.Text)

	var reply string

	switch state {

	//case StateMain:
	//	reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
	//	messageConfig := h.setInlineKeyboard(message.Chat.ID, reply, []string{"перейти к списку монет", "перейти в избранное"})
	//	return messageConfig
	//
	//case StateCoinsList:
	//
	//	reply = "список всех монет:"
	//
	//	coinSymbols, symbolsGettingError := h.service.GetCoinsSymbols()
	//	if symbolsGettingError != nil {
	//		log.Println(symbolsGettingError)
	//	}
	//
	//	messageConfig := h.setCoinsKeyboard(message.Chat.ID, reply, coinSymbols)
	//	return messageConfig
	//
	//case StateCoinInfo:
	//
	//	coinSymbol := message.Text
	//	coinInfo, gettingSymbolError := h.service.GetCoinInfo(coinSymbol)
	//	if gettingSymbolError != nil {
	//		log.Println(gettingSymbolError)
	//		return tgbotapi.NewMessage(message.Chat.ID, "ошибка получения данных")
	//	}
	//
	//	reply := fmt.Sprintf("информация о криптовалюте %s:\nназвание: %s\nцена: %s$\nизменение цены за 1 час: %s\nизменение цены за 24 часа: %s\nизменение цены за 7 дней: %s\n", coinSymbol, coinInfo.Name, coinInfo.PriceUsd, coinInfo.PercentChange1H, coinInfo.PercentChange24H, coinInfo.PercentChange7D)
	//
	//	messageConfig := h.setInlineKeyboard(message.Chat.ID, reply, []string{"назад"})
	//
	//	return messageConfig
	case StateSetAlert:
		reply = fmt.Sprintf("Монета теперь отслеживается")

		id := fmt.Sprint(message.From.ID)
		h.service.AddPriceAlert(id, *currentCoin, message.Text)
		return h.setInlineKeyboard(message.Chat.ID, reply, []string{"далее"})

	default:
		return tgbotapi.NewMessage(message.Chat.ID, "")
	}

}

func (h *Handler) handleCallbackQuery(callback *tgbotapi.CallbackQuery, currentCoin *string) tgbotapi.MessageConfig {
	state := h.stateMachine.SetState(callback.Data)
	fmt.Println(callback.Message)

	var reply string
	switch state {
	case StateMain:
		reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
		messageConfig := h.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"перейти к списку монет", "избранное"})
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

		coinForecast, forecastError := h.service.GetCoinForecast(coinSymbol)
		if forecastError != nil {
			log.Println(forecastError)
			return tgbotapi.NewMessage(callback.Message.Chat.ID, "ошибка получения данных")
		}

		pr := "%"
		reply := fmt.Sprintf("информация о криптовалюте %s:\nназвание: %s\nцена: %s$\nизменение цены за 1 час: %s %s\nизменение цены за 24 часа: %s %s\nизменение цены за 7 дней: %s %s\nпрогноз: %s $", coinSymbol, coinInfo.Name, coinInfo.PriceUsd, coinInfo.PercentChange1H, pr, coinInfo.PercentChange24H, pr, coinInfo.PercentChange7D, pr, coinForecast)
		var buttonsMap map[string]string
		buttonsMap = make(map[string]string)
		buttonsMap["добавить в избранное"] = "@" + callback.Data
		buttonsMap["установить alert"] = callback.Data
		buttonsMap["назад"] = "назад"

		messageConfig := h.setCoinInfoKeyboard(callback.Message.Chat.ID, reply, buttonsMap)

		return messageConfig

	case StateCollection:
		reply = "ваши избранные монеты:\n"

		id := fmt.Sprint(callback.Message.From.ID)
		collection := h.service.GetCollection(id)
		var coinSymbols []string
		for _, coin := range collection {
			coinSymbols = append(coinSymbols, coin.Symbol)
		}

		messageConfig := h.setCoinsKeyboard(callback.Message.Chat.ID, reply, coinSymbols)
		return messageConfig

	case StateSetCollection:
		reply = fmt.Sprintf("Монета добавлена в избранное")

		id := fmt.Sprint(callback.Message.From.ID)
		h.service.SetCollectionItem(id, callback.Data[1:])
		return h.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"далее"})

	case StateRemoveCollection:
		reply = fmt.Sprintf("Монета удалена из избранного")

		id := fmt.Sprint(callback.Message.From.ID)
		h.service.SetCollectionItem(id, callback.Data[1:])
		return h.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"далее"})

	case StateAlert:
		currentCoin := callback.Data
		reply = fmt.Sprintf("Введите значение %s:", currentCoin)
		return tgbotapi.NewMessage(callback.Message.Chat.ID, reply)
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

func (h *Handler) setCoinInfoKeyboard(chatId int64, text string, buttons map[string]string) tgbotapi.MessageConfig {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	//for key, value := range buttons {
	//	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(key, value)))
	//}

	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("добавить в избранное", buttons["добавить в избранное"])))
	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("установить alert", buttons["установить alert"])))
	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "назад")))

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	message := tgbotapi.NewMessage(chatId, text)
	message.ReplyMarkup = inlineKeyboard
	return message
}
