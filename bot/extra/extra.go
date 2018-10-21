package extra

import (
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/json-iterator/go"
)

func UpdateMsgPhoto(
	bot *tgbotapi.BotAPI, chatID int64, messageID int,
	caption, parseMode string,
	markup tgbotapi.InlineKeyboardMarkup, file interface{},
) error {
	media, err := jsoniter.MarshalToString(struct {
		Type      string `json:"type"`
		Media     string `json:"media"`
		Caption   string `json:"caption"`
		ParseMode string `json:"parse_mode"`
	}{
		Type:      "photo",
		Media:     "attach://photo",
		Caption:   caption,
		ParseMode: parseMode,
	})
	if err != nil {
		return err
	}
	replyMarkup, err := jsoniter.MarshalToString(markup)
	if err != nil {
		return err
	}

	_, err = bot.UploadFile("editMessageMedia", map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"message_id":   strconv.Itoa(messageID),
		"media":        media,
		"reply_markup": replyMarkup,
	}, "photo", file)
	return err
}
