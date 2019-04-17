package extra

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/internal"
	"github.com/jqs7/zwei/model"
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
	mediaReq, err := json.Marshal(struct {
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
	replyMarkup, err := json.Marshal(markup)
	if err != nil {
		return nil, err
	}

	reqParam := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"message_id":   strconv.Itoa(messageID),
		"media":        string(mediaReq),
		"reply_markup": string(replyMarkup),
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
		return message, json.Unmarshal(resp.Result, message)
	}
	resp, err := bot.UploadFile("editMessageMedia", reqParam, "photo", file)
	if err != nil {
		return nil, err
	}
	message := &tgbotapi.Message{}
	return message, json.Unmarshal(resp.Result, message)
}

func KickAndDelCaptcha(bot *tgbotapi.BotAPI, blackList model.BlackList, banUntil int64) {
	internal.JustLogErr(bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId)))
	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID: blackList.GroupId,
		UserID: blackList.UserId,
	})
	if err == nil && (member.HasLeft() || member.WasKicked()) {
		return
	}
	internal.JustLogErr(bot.KickChatMember(tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: blackList.GroupId,
			UserID: blackList.UserId,
		},
		UntilDate: banUntil,
	}))
}
