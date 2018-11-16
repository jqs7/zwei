package extra

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/model"
	"github.com/json-iterator/go"
)

func UpdateMsgPhoto(
	bot *tgbotapi.BotAPI, chatID int64, messageID int,
	caption, parseMode string,
	markup tgbotapi.InlineKeyboardMarkup, file interface{},
) (*tgbotapi.Message, error) {
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
		return nil, err
	}
	replyMarkup, err := jsoniter.MarshalToString(markup)
	if err != nil {
		return nil, err
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
		resp, err := bot.MakeRequest("editMessageMedia", values)
		if err != nil {
			return nil, err
		}
		message := &tgbotapi.Message{}
		json.Unmarshal(resp.Result, message)
		return message, nil
	}
	resp, err := bot.UploadFile("editMessageMedia", reqParam, "photo", file)
	if err != nil {
		return nil, err
	}
	message := &tgbotapi.Message{}
	json.Unmarshal(resp.Result, message)
	return message, nil
}

func KickAndDelCaptcha(bot *tgbotapi.BotAPI, blackList model.BlackList) {
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId))
	bot.KickChatMember(tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: blackList.GroupId,
			UserID: blackList.UserId,
		},
		UntilDate: time.Now().Add(time.Minute).Unix(),
	})
}
