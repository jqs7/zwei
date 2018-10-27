package model

import (
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	CallbackTypeRefresh     = "Refresh"
	CallbackTypePassThrough = "PassThrough"
)

const (
	UserLinkTemplate = "[%s](tg://user?id=%d)"
	EnterRoomMsg     = `%s 你好，欢迎加入 %s，本群已启用新成员验证模式，请发送图片验证码内容。
在验证通过之前，你所发送的所有消息都将会被删除。
本消息将在 %d 秒后失效，届时若未通过验证，你将被移出群组，且一分钟之内无法再加入本群。`
	HelpMsg = `欢迎使用进群验证码机器人
本机器人使用姿势：
将本机器人加入需要启用验证的群组，设置为管理员，并授予 Delete messages，Ban users 权限即可
本项目开源于：https://github.com/jqs7/zwei`
)

const (
	DefaultCaptchaExpire   = 300 * time.Second
	DefaultRefreshDuration = 75 * time.Second
)

var InlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("刷新", CallbackTypeRefresh),
		tgbotapi.NewInlineKeyboardButtonData("通过验证[管理员]", CallbackTypePassThrough),
	},
)
