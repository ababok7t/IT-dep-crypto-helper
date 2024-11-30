package services

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (s *Service) setInlineKeyboard(chatId int64, text string, buttons []string) tgbotapi.MessageConfig {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	for _, button := range buttons {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(button, button)))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	message := tgbotapi.NewMessage(chatId, text)
	message.ReplyMarkup = inlineKeyboard
	return message
}

func (s *Service) setCoinsKeyboard(chatId int64, text string, buttons []string) tgbotapi.MessageConfig {
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

func (s *Service) setCoinInfoKeyboard(chatId int64, text string, buttons map[string]string) tgbotapi.MessageConfig {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	//for key, value := range buttons {
	//	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(key, value)))
	//}

	_, err := buttons["добавить в избранное"]
	if err {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("добавить в избранное", buttons["добавить в избранное"])))
	} else {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("удалить из избранного", buttons["удалить из избранного"])))
	}

	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("установить alert", buttons["установить alert"])))
	inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "назад")))

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
	message := tgbotapi.NewMessage(chatId, text)
	message.ReplyMarkup = inlineKeyboard
	return message
}
