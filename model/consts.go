package model

import (
	"sort"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	CallbackTypeRefresh     = "Refresh"
	CallbackTypePassThrough = "PassThrough"
	CallbackTypeKick        = "Kick"

	CallbackTypeDonateWX     = "DonateWX"
	CallbackTypeDonateAlipay = "DonateAlipay"
)

const (
	UserLinkTemplate = "[%s](tg://user?id=%d)"
	EnterRoomMsg     = `%s 你好，欢迎加入 %s，本群已启用新成员验证模式，请发送以上 *【四字】* 验证码内容。
在验证通过之前，你所发送的所有消息都将会被删除。
本消息将在 %d 秒后失效，届时若未通过验证，你将被移出群组，且一分钟之内无法再加入本群。`
	HelpMsg = `欢迎使用进群验证码机器人
本机器人使用姿势：
将本机器人加入需要启用验证的群组，设置为管理员，并授予 Delete messages，Ban users 权限即可
本项目开源于：https://github.com/jqs7/zwei
若本项目对你有所帮助，可点击 /donate 为本项目捐款`
	DonateMsg = `所捐款项将用于：
1. 作者的续命咖啡 ☕️
2. 支付服务器等设施费用`
)

const (
	DefaultCaptchaExpire   = 300 * time.Second
	DefaultRefreshDuration = 75 * time.Second
)

var InlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("刷新验证码", CallbackTypeRefresh),
		tgbotapi.NewInlineKeyboardButtonData("通过验证[管理员]", CallbackTypePassThrough),
	},
	[]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("踢出群组[管理员]", CallbackTypeKick),
	},
)

type DonateKV struct {
	Key    string
	FileID string
	URL    string
}

var Donates = map[string]DonateKV{
	CallbackTypeDonateWX: {
		Key: "微信",
		URL: "wxp://f2f0OWfabxt-G2eVGJuF9psyiEvqiL3u3gxB",
	},
	CallbackTypeDonateAlipay: {
		Key: "支付宝",
		URL: "https://qr.alipay.com/fkx00824kg0dc3tf1sf4c2e",
	},
}

type InlineKeyboardButtons []tgbotapi.InlineKeyboardButton

func (iBtn InlineKeyboardButtons) Len() int {
	return len(iBtn)
}

func (iBtn InlineKeyboardButtons) Less(i, j int) bool {
	return iBtn[i].Text < iBtn[j].Text
}

func (iBtn InlineKeyboardButtons) Swap(i, j int) {
	iBtn[i], iBtn[j] = iBtn[j], iBtn[i]
}

func DonatesKeyboard(donateType string) tgbotapi.InlineKeyboardMarkup {
	var buttons InlineKeyboardButtons
	for k, v := range Donates {
		if k == donateType {
			v.Key = v.Key + "❤️"
		}
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v.Key, k))
	}
	sort.Sort(buttons)
	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}
