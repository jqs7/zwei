package biz

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	CallbackTypeRefresh     = "Refresh"
	CallbackTypePassThrough = "PassThrough"
)

var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("刷新", CallbackTypeRefresh),
		tgbotapi.NewInlineKeyboardButtonData("通过验证[管理员]", CallbackTypePassThrough),
	},
)

var UserLinkTemplate = "[%s](tg://user?id=%d)"
var EnterRoomMsg = "%s 你好，欢迎加入 %s，本群已启用新成员验证模式，请发送图片验证码内容。\n" +
	"在验证通过之前，你所发送的所有消息都将会被删除。"
