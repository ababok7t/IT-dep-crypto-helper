package services

import (
	"crypto-helper/internal/domain"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (s *Service) HandleMessage(message tgbotapi.Message, currentCoin *string) tgbotapi.MessageConfig {
	state := s.StateMachine.SetState(message.Text)

	var reply string

	switch state {
	case StateMain:
		reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
		messageConfig := s.setInlineKeyboard(message.Chat.ID, reply, []string{"перейти к списку монет", "избранное"})
		return messageConfig

	case StateSetAlert:
		reply = fmt.Sprintf("Монета теперь отслеживается")

		id := fmt.Sprint(message.From.ID)
		s.AddPriceAlert(id, *currentCoin, message.Text)
		return s.setInlineKeyboard(message.Chat.ID, reply, []string{"далее"})

	default:
		return tgbotapi.NewMessage(message.Chat.ID, "")
	}

}

func (s *Service) HandleCallbackQuery(callback *tgbotapi.CallbackQuery, currentCoin *string) tgbotapi.MessageConfig {
	state := s.StateMachine.SetState(callback.Data)
	fmt.Println(callback.Message)

	var reply string
	switch state {
	case StateMain:
		reply = "Crypto Helper Bot\nВаш помощник в мире криптовалютных инвестиций"
		messageConfig := s.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"перейти к списку монет", "избранное"})
		return messageConfig

	case StateCoinsList:

		reply = "список всех монет:"

		coinSymbols, symbolsGettingError := s.GetCoinsSymbols()
		if symbolsGettingError != nil {
			log.Println(symbolsGettingError)
		}

		messageConfig := s.setCoinsKeyboard(callback.Message.Chat.ID, reply, coinSymbols)
		return messageConfig

	case StateCoinInfo:

		coinSymbol := callback.Data
		coinInfo, gettingSymbolError := s.GetCoinInfo(coinSymbol)
		if gettingSymbolError != nil {
			log.Println(gettingSymbolError)
			return tgbotapi.NewMessage(callback.Message.Chat.ID, "ошибка получения данных")
		}

		coinForecast, forecastError := s.GetCoinForecast(coinSymbol)
		if forecastError != nil {
			log.Println(forecastError)
			return tgbotapi.NewMessage(callback.Message.Chat.ID, "ошибка получения данных")
		}

		pr := "%"
		reply := domain.MakeCoinReply(coinSymbol, coinInfo, coinForecast, pr)

		var buttonsMap map[string]string
		buttonsMap = make(map[string]string)
		id := fmt.Sprint(callback.Message.From.ID)
		flag := 0
		for _, coin := range s.GetCollection(id) {
			if callback.Data == coin.Symbol {
				flag = 1
			}
		}

		if flag == 0 {
			buttonsMap["добавить в избранное"] = "@" + callback.Data
		}

		if flag == 1 {
			buttonsMap["удалить из избранного"] = "-" + callback.Data
		}

		buttonsMap["установить alert"] = callback.Data
		buttonsMap["назад"] = "назад"

		messageConfig := s.setCoinInfoKeyboard(callback.Message.Chat.ID, reply, buttonsMap)

		return messageConfig

	case StateCollection:
		reply = "ваши избранные монеты:\n"

		id := fmt.Sprint(callback.Message.From.ID)
		collection := s.GetCollection(id)
		var coinSymbols []string
		for _, coin := range collection {
			coinSymbols = append(coinSymbols, coin.Symbol)
		}

		messageConfig := s.setCoinsKeyboard(callback.Message.Chat.ID, reply, coinSymbols)
		return messageConfig

	case StateSetCollection:
		reply = fmt.Sprintf("Монета добавлена в избранное")

		id := fmt.Sprint(callback.Message.From.ID)
		s.SetCollectionItem(id, callback.Data[1:])
		return s.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"далее"})

	case StateRemoveCollection:
		reply = fmt.Sprintf("Монета удалена из избранного")

		id := fmt.Sprint(callback.Message.From.ID)
		s.RemoveCollectionItem(id, callback.Data[1:])
		return s.setInlineKeyboard(callback.Message.Chat.ID, reply, []string{"далее"})

	case StateAlert:
		currentCoin := callback.Data
		reply = fmt.Sprintf("Введите значение %s:", currentCoin)
		return tgbotapi.NewMessage(callback.Message.Chat.ID, reply)
	default:
		return tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	}

}
