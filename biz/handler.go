package biz

import (
	"bytes"
	"fmt"
	"html"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hanguofeng/gocaptcha"
	"github.com/jqs7/zwei/bot/extra"
	"github.com/jqs7/zwei/db"
	"github.com/jqs7/zwei/env"
	"github.com/jqs7/zwei/internal"
	"github.com/jqs7/zwei/model"
	"github.com/jqs7/zwei/scheduler"
	"github.com/skip2/go-qrcode"
)

type Handler struct {
	*db.DB
	*gocaptcha.ImageConfig
	*gocaptcha.ImageFilterManager
	IdiomCount int
}

func NewHandler(db *db.DB, cfg env.Specification) Handler {
	idiomCount, err := db.PgDB.Model(new(model.Idiom)).Count()
	if err != nil {
		log.Fatalln(err)
	}

	filterConfig := new(gocaptcha.FilterConfig)
	filterConfig.Init()
	filterConfig.Filters = []string{
		gocaptcha.IMAGE_FILTER_NOISE_LINE,
		gocaptcha.IMAGE_FILTER_NOISE_POINT,
		gocaptcha.IMAGE_FILTER_STRIKE,
	}
	for _, v := range filterConfig.Filters {
		filterConfigGroup := new(gocaptcha.FilterConfigGroup)
		filterConfigGroup.Init()
		filterConfigGroup.SetItem("Num", strconv.Itoa(cfg.CaptchaNoise))
		filterConfig.SetGroup(v, filterConfigGroup)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fontPath := filepath.Join(pwd, "fonts")
	return Handler{
		DB: db,
		ImageConfig: &gocaptcha.ImageConfig{
			Width:    320,
			Height:   100,
			FontSize: 80,
			FontFiles: []string{
				filepath.Join(fontPath, "STFANGSO.ttf"),
				filepath.Join(fontPath, "STHEITI.ttf"),
				filepath.Join(fontPath, "STXIHEI.ttf"),
			},
		},
		IdiomCount:         idiomCount,
		ImageFilterManager: gocaptcha.CreateImageFilterManagerByConfig(filterConfig),
	}
}

func (Handler) BotEnterGroup(*tgbotapi.BotAPI, *tgbotapi.Chat) error {
	return nil
}

func (h Handler) OnPrivateCommand(bot *tgbotapi.BotAPI, msg tgbotapi.Message, command string, args ...string) error {
	switch command {
	case "help", "start":
		return h.sendHelpMsg(bot, msg.Chat.ID)
	case "donate":
		return h.sendDonate(bot, msg.Chat.ID, model.CallbackTypeDonateWX)
	}
	return nil
}

func (h Handler) sendDonate(bot *tgbotapi.BotAPI, chatID int64, donateType string) error {
	var donateImg tgbotapi.PhotoConfig
	if model.Donates[donateType].FileID != "" {
		donateImg = tgbotapi.NewPhotoShare(chatID, model.Donates[donateType].FileID)
	} else {
		b, err := qrcode.Encode(model.Donates[donateType].URL, qrcode.Medium, 256)
		if err != nil {
			return err
		}
		donateImg = tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "file", Bytes: b})
	}
	donateImg.ReplyMarkup = model.DonatesKeyboard(donateType)
	donateImg.Caption = model.DonateMsg
	msg, err := bot.Send(donateImg)
	if err != nil {
		return err
	}
	donate := model.Donates[donateType]
	donate.FileID = (*msg.Photo)[len(*msg.Photo)-1].FileID
	model.Donates[donateType] = donate
	return err
}

func (h Handler) sendHelpMsg(bot *tgbotapi.BotAPI, toUserID int64) error {
	msg := tgbotapi.NewMessage(toUserID, model.HelpMsg)
	msg.DisableNotification = true
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	return err
}

func (h Handler) NewMemberInGroup(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, user tgbotapi.User) error {
	blackList := &model.BlackList{
		GroupId: chat.ID,
		UserId:  user.ID,
	}
	err := h.PgDB.Insert(blackList)
	if err != nil {
		return err
	}

	idiom, err := h.GetRandomIdiom()
	if err != nil {
		return err
	}
	return h.sendCaptcha(bot, chat, user, blackList, idiom)
}

func (h Handler) OnMemberLeftGroup(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, user tgbotapi.User) error {
	var blackLists []model.BlackList
	if err := h.PgDB.Model(&blackLists).
		Where("group_id = ?", chat.ID).
		Where("user_id = ?", user.ID).
		Select(); err != nil && err != pg.ErrNoRows {
		return err
	}
	if len(blackLists) == 0 {
		return nil
	}
	for _, blackList := range blackLists {
		internal.JustLogErr(bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId)))
		internal.JustLogErr(scheduler.UpdateMsgExpireTaskDone(h.DB, blackList.Id))
	}
	_, err := h.PgDB.Model(&model.BlackList{}).
		Where("group_id = ?", chat.ID).
		Where("user_id = ?", user.ID).
		Delete()
	return err
}

func (h Handler) sendCaptcha(bot *tgbotapi.BotAPI,
	chat *tgbotapi.Chat, user tgbotapi.User,
	blackList *model.BlackList, idiom *model.Idiom,
) error {
	captchaMsg := tgbotapi.NewPhotoUpload(chat.ID, tgbotapi.FileBytes{
		Name:  strconv.Itoa(user.ID),
		Bytes: idiom.CaptchaImg,
	})
	fullName := getFullName(user.FirstName, user.LastName)
	userLink := fmt.Sprintf(model.UserLinkTemplate, user.ID, html.EscapeString(fullName))
	captchaMsg.Caption = fmt.Sprintf(model.EnterRoomMsg, userLink, chat.Title, model.DefaultCaptchaExpire/time.Second)
	captchaMsg.ParseMode = tgbotapi.ModeHTML
	captchaMsg.ReplyMarkup = model.InlineKeyboard
	captchaMsg.DisableNotification = true
	photoMsg, err := bot.Send(captchaMsg)
	if err != nil {
		return err
	}
	if err := scheduler.AddUpdateMsgExpireTask(
		h.DB, blackList.Id, chat.ID, photoMsg.MessageID,
	); err != nil {
		return err
	}
	blackList.IdiomId = idiom.Id
	blackList.CaptchaMsgId = photoMsg.MessageID
	blackList.UserLink = userLink
	blackList.ExpireAt = time.Now().Add(model.DefaultCaptchaExpire)
	if _, err := h.PgDB.Model(blackList).
		Column("idiom_id", "captcha_msg_id", "user_link", "expire_at").
		WherePK().Update(); err != nil {
		return err
	}
	return nil
}

func (h Handler) OnGroupMsg(bot *tgbotapi.BotAPI, msg tgbotapi.Message) error {
	blackList := &model.BlackList{}
	err := h.PgDB.Model(blackList).
		Column("black_list.*", "Idiom").
		Where("group_id = ?", msg.Chat.ID).
		Where("user_id = ?", msg.From.ID).
		Last()
	if err == pg.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	internal.JustLogErr(h.deleteMsg(bot, msg.Chat.ID, msg.MessageID))
	if msg.Text != blackList.Idiom.Word {
		return nil
	}
	return h.validateOK(bot, blackList)
}

func (h Handler) deleteMsg(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	_, err := bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, messageID))
	return err
}

func (h Handler) OnCallbackQuery(bot *tgbotapi.BotAPI, query tgbotapi.CallbackQuery) error {
	switch query.Data {
	case model.CallbackTypeRefresh:
		blackList := &model.BlackList{}
		if err := h.PgDB.Model(blackList).
			Where("group_id = ?", query.Message.Chat.ID).
			Where("captcha_msg_id = ?", query.Message.MessageID).
			First(); err != nil {
			return err
		}
		return h.refresh(bot, blackList, query)
	case model.CallbackTypePassThrough:
		blackList := &model.BlackList{}
		if err := h.PgDB.Model(blackList).
			Where("group_id = ?", query.Message.Chat.ID).
			Where("captcha_msg_id = ?", query.Message.MessageID).
			First(); err != nil {
			return err
		}
		return h.passThrough(bot, blackList, query)
	case model.CallbackTypeKick:
		member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
			ChatID:             query.Message.Chat.ChatConfig().ChatID,
			SuperGroupUsername: query.Message.Chat.ChatConfig().SuperGroupUsername,
			UserID:             query.From.ID,
		})
		if err != nil {
			return err
		}
		if !member.IsCreator() && !member.IsAdministrator() {
			return h.answerCallbackQuery(bot, query, "无权限")
		}
		blackList := &model.BlackList{}
		if err := h.PgDB.Model(blackList).
			Where("group_id = ?", query.Message.Chat.ID).
			Where("captcha_msg_id = ?", query.Message.MessageID).
			First(); err != nil {
			return err
		}
		_, err = h.PgDB.Model(blackList).
			Where("group_id = ?group_id").
			Where("user_id = ?user_id").
			Delete()
		if err != nil {
			return err
		}
		internal.JustLogErr(scheduler.UpdateMsgExpireTaskDone(h.DB, blackList.Id))
		extra.KickAndDelCaptcha(bot, *blackList, 0)
	case model.CallbackTypeDonateWX, model.CallbackTypeDonateAlipay:
		if model.Donates[query.Data].FileID != "" {
			_, err := extra.UpdateMsgPhoto(bot, query.Message.Chat.ID, query.Message.MessageID, query.Message.Caption,
				"", model.DonatesKeyboard(query.Data), model.Donates[query.Data].FileID)
			if err != nil {
				return err
			}
		} else {
			b, err := qrcode.Encode(model.Donates[query.Data].URL, qrcode.Medium, 256)
			if err != nil {
				return err
			}
			msg, err := extra.UpdateMsgPhoto(bot, query.Message.Chat.ID, query.Message.MessageID, query.Message.Caption,
				"", model.DonatesKeyboard(query.Data), tgbotapi.FileBytes{Name: "file", Bytes: b})
			if err != nil {
				return err
			}
			donate := model.Donates[query.Data]
			donate.FileID = (*msg.Photo)[len(*msg.Photo)-1].FileID
			model.Donates[query.Data] = donate
		}
		_, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, ""))
		return err
	}
	return nil
}

func (h Handler) refresh(bot *tgbotapi.BotAPI, blackList *model.BlackList, query tgbotapi.CallbackQuery) error {
	if blackList.ExpireAt.Before(time.Now()) {
		return h.answerCallbackQuery(bot, query, "已过期")
	}
	if query.From.ID != blackList.UserId {
		return h.answerCallbackQuery(bot, query, "无权限")
	}
	idiom, err := h.GetRandomIdiom()
	if err != nil {
		internal.JustLogErr(h.answerCallbackQuery(bot, query, "刷新失败"))
		return err
	}
	blackList.IdiomId = idiom.Id
	if _, err := h.PgDB.Model(blackList).
		Column("idiom_id").
		WherePK().Update(); err != nil {
		internal.JustLogErr(h.answerCallbackQuery(bot, query, "刷新失败"))
		return err
	}
	if _, err := extra.UpdateMsgPhoto(
		bot, query.Message.Chat.ID, query.Message.MessageID,
		fmt.Sprintf(model.EnterRoomMsg, blackList.UserLink, query.Message.Chat.Title,
			time.Until(blackList.ExpireAt)/time.Second),
		tgbotapi.ModeHTML, model.InlineKeyboard, tgbotapi.FileBytes{
			Name:  strconv.Itoa(query.From.ID),
			Bytes: idiom.CaptchaImg,
		},
	); err != nil {
		internal.JustLogErr(h.answerCallbackQuery(bot, query, "刷新失败"))
		return err
	}
	return h.answerCallbackQuery(bot, query, "刷新成功")
}

func (h Handler) answerCallbackQuery(bot *tgbotapi.BotAPI, query tgbotapi.CallbackQuery, answer string) error {
	_, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, answer))
	return err
}

func (h Handler) passThrough(bot *tgbotapi.BotAPI, blackList *model.BlackList, query tgbotapi.CallbackQuery) error {
	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID:             query.Message.Chat.ChatConfig().ChatID,
		SuperGroupUsername: query.Message.Chat.ChatConfig().SuperGroupUsername,
		UserID:             query.From.ID,
	})
	if err != nil {
		return err
	}
	if !member.IsCreator() && !member.IsAdministrator() {
		return h.answerCallbackQuery(bot, query, "无权限")
	}
	return h.validateOK(bot, blackList)
}

func (h Handler) validateOK(bot *tgbotapi.BotAPI, blackList *model.BlackList) error {
	_, err := h.PgDB.Model(blackList).
		Where("group_id = ?group_id").
		Where("user_id = ?user_id").
		Delete()
	if err != nil {
		return err
	}
	internal.JustLogErr(bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId)))
	passMsgConfig := tgbotapi.NewMessage(blackList.GroupId, fmt.Sprintf(
		"%s 恭喜，你已验证通过", blackList.UserLink))
	passMsgConfig.DisableNotification = true
	passMsgConfig.ParseMode = tgbotapi.ModeHTML
	passMsg, err := bot.Send(passMsgConfig)
	if err != nil {
		return err
	}
	internal.JustLogErr(scheduler.UpdateMsgExpireTaskDone(h.DB, blackList.Id))
	return scheduler.AddDelMsgTask(h.DB, blackList.GroupId, passMsg.MessageID)
}

func (h Handler) GetRandomIdiom() (*model.Idiom, error) {
	idiom := &model.Idiom{}
	rand.Seed(time.Now().UnixNano())
	randOffset := rand.Intn(h.IdiomCount)
	if err := h.PgDB.Model(idiom).
		Offset(randOffset).Limit(1).Select(); err != nil {
		return nil, err
	}
	cImg := gocaptcha.CreateCImage(h.ImageConfig)
	cImg.DrawString(idiom.Word)
	for _, f := range h.GetFilters() {
		f.Proc(cImg)
	}
	buffer := bytes.NewBuffer([]byte{})
	if err := png.Encode(buffer, cImg); err != nil {
		return nil, err
	}
	idiom.CaptchaImg = buffer.Bytes()
	return idiom, nil
}

func getFullName(firstName, lastName string) string {
	fullName := strings.TrimSpace(firstName + " " + lastName)
	if len([]rune(fullName)) > 10 {
		fullName = string([]rune(fullName)[:10]) + "..."
	}
	return fullName
}
