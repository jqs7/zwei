package tg

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/biz"
)

type Bot struct {
	*tgbotapi.BotAPI
	myInfo tgbotapi.User
	biz.IBiz
	updates tgbotapi.UpdatesChannel
}

func (b Bot) processUpdate(update tgbotapi.Update) error {
	if update.Message != nil {
		return b.onMessage(*update.Message)
	}
	if update.CallbackQuery != nil {
		return b.OnCallbackQuery(b.BotAPI, *update.CallbackQuery)
	}
	return nil
}

func (b Bot) onMessage(msg tgbotapi.Message) error {
	if msg.GroupChatCreated {
		return b.BotEnterGroup(b.BotAPI, msg.Chat)
	}
	if msg.NewChatMembers != nil {
		return b.onNewChatMembers(msg)
	}
	if msg.LeftChatMember != nil {
		return b.onLeftChatMember(msg)
	}
	if msg.Chat.IsGroup() || msg.Chat.IsSuperGroup() {
		return b.OnGroupMsg(b.BotAPI, msg)
	}
	if msg.IsCommand() {
		return b.OnPrivateCommand(b.BotAPI, msg,
			msg.Command(), strings.Split(msg.CommandArguments(), " ")...,
		)
	}
	return nil
}

func (b Bot) onNewChatMembers(msg tgbotapi.Message) error {
	b.DeleteMessage(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
	for _, v := range *msg.NewChatMembers {
		if v.ID == b.myInfo.ID {
			return b.BotEnterGroup(b.BotAPI, msg.Chat)
		}
		if v.IsBot {
			continue
		}
		if err := b.NewMemberInGroup(b.BotAPI, msg.Chat, v); err != nil {
			return err
		}
	}
	return nil
}

func (b Bot) onLeftChatMember(msg tgbotapi.Message) error {
	b.DeleteMessage(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
	return b.OnMemberLeftGroup(b.BotAPI, msg.Chat, *msg.LeftChatMember)
}
