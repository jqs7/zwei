package biz

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type IBiz interface {
	BotEnterGroup(*tgbotapi.BotAPI, *tgbotapi.Chat) error
	NewMemberInGroup(*tgbotapi.BotAPI, *tgbotapi.Chat, tgbotapi.User) error
	OnGroupMsg(*tgbotapi.BotAPI, tgbotapi.Message) error
	OnCallbackQuery(*tgbotapi.BotAPI, tgbotapi.CallbackQuery) error
	OnMemberLeftGroup(*tgbotapi.BotAPI, *tgbotapi.Chat, tgbotapi.User) error
}
