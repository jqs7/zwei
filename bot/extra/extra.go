package extra

import (
	"net/url"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/json-iterator/go"
)

func UpdateMsgPhoto(
	bot *tgbotapi.BotAPI, chatID int64, messageID int,
	caption, parseMode string,
	markup tgbotapi.InlineKeyboardMarkup, file interface{},
) error {
	media := "attach://photo"
	fileID, withFileID := file.(string)
	if withFileID {
		media = fileID
	}
	media, err := jsoniter.MarshalToString(struct {
		Type      string `json:"type"`
		Media     string `json:"media"`
		Caption   string `json:"caption"`
		ParseMode string `json:"parse_mode"`
	}{
		Type:      "photo",
		Media:     media,
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

	reqParam := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"message_id":   strconv.Itoa(messageID),
		"media":        media,
		"reply_markup": replyMarkup,
	}

	if withFileID {
		values := url.Values{}
		for k, v := range reqParam {
			values.Set(k, v)
		}
		_, err := bot.MakeRequest("editMessageMedia", values)
		return err
	}
	_, err = bot.UploadFile("editMessageMedia", reqParam, "photo", file)
	return err
}
