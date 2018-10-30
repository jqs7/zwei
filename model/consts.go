package model

import (
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	CallbackTypeRefresh     = "Refresh"
	CallbackTypePassThrough = "PassThrough"

	CallbackTypeDonateWX     = "DonateWX"
	CallbackTypeDonateAlipay = "DonateAlipay"
)

const (
	UserLinkTemplate = "[%s](tg://user?id=%d)"
	EnterRoomMsg     = `%s 你好，欢迎加入 %s，本群已启用新成员验证模式，请发送以上 *四字* 验证码内容。
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

type DonateKV struct {
	Key string
	ID  string
}

var Donates = map[string]DonateKV{
	CallbackTypeDonateWX: {
		Key: "微信",
		ID:  "AgADBQADI6gxG_OVyVZPEk79HZiSzz9h2zIABA2NBWY3mfZlkOwAAgI",
	},
	CallbackTypeDonateAlipay: {
		Key: "支付宝",
		ID:  "AgADBQADIqgxG_OVyVZoTf04FO5TuWhm2zIABBeF-IkG4kvIBucAAgI",
	},
}

func DonatesKeyboard(donateType string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	for k, v := range Donates {
		if k == donateType {
			v.Key = v.Key + "❤️"
		}
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Key, k))
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}
