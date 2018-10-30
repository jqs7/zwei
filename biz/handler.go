package biz

import (
	"bytes"
	"fmt"
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
	"github.com/jqs7/zwei/model"
	"github.com/jqs7/zwei/scheduler"
)

type Handler struct {
	*gocaptcha.ImageConfig
	*gocaptcha.ImageFilterManager
	IdiomCount int
}

func NewHandler(noiseNum int) Handler {
	idiomCount, err := db.Instance().Model(new(model.Idiom)).Count()
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
		filterConfigGroup.SetItem("Num", strconv.Itoa(noiseNum))
		filterConfig.SetGroup(v, filterConfigGroup)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fontPath := filepath.Join(pwd, "fonts")
	return Handler{
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
	donateImg := tgbotapi.NewPhotoShare(chatID, model.Donates[donateType].ID)
	donateImg.ReplyMarkup = model.DonatesKeyboard(donateType)
	bot.Send(donateImg)
	return nil
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
	err := db.Instance().Insert(blackList)
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
	if err := db.Instance().Model(&blackLists).
		Where("group_id = ?", chat.ID).
		Where("user_id = ?", user.ID).
		Select(); err != nil && err != pg.ErrNoRows {
		return err
	}
	if len(blackLists) == 0 {
		return nil
	}
	for _, blackList := range blackLists {
		bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId))
		scheduler.UpdateMsgExpireTaskDone(db.Instance(), blackList.Id)
	}
	_, err := db.Instance().Model(&model.BlackList{}).
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
	userLink := fmt.Sprintf(model.UserLinkTemplate, fullName, user.ID)
	captchaMsg.Caption = fmt.Sprintf(model.EnterRoomMsg, userLink, chat.Title, model.DefaultCaptchaExpire/time.Second)
	captchaMsg.ParseMode = tgbotapi.ModeMarkdown
	captchaMsg.ReplyMarkup = model.InlineKeyboard
	captchaMsg.DisableNotification = true
	photoMsg, err := bot.Send(captchaMsg)
	if err != nil {
		return err
	}
	if err := scheduler.AddUpdateMsgExpireTask(
		db.Instance(), blackList.Id, chat.ID, photoMsg.MessageID,
	); err != nil {
		return err
	}
	blackList.IdiomId = idiom.Id
	blackList.CaptchaMsgId = photoMsg.MessageID
	blackList.UserLink = userLink
	blackList.ExpireAt = time.Now().Add(model.DefaultCaptchaExpire)
	if _, err := db.Instance().Model(blackList).
		Column("idiom_id", "captcha_msg_id", "user_link", "expire_at").
		WherePK().Update(); err != nil {
		return err
	}
	return nil
}

func (h Handler) OnGroupMsg(bot *tgbotapi.BotAPI, msg tgbotapi.Message) error {
	blackList := &model.BlackList{}
	err := db.Instance().Model(blackList).
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
	h.deleteMsg(bot, msg.Chat.ID, msg.MessageID)
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
		err := db.Instance().Model(blackList).
			Where("group_id = ?", query.Message.Chat.ID).
			Where("captcha_msg_id = ?", query.Message.MessageID).
			First()
		if err != nil {
			_, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, err.Error()))
			return err
		}
		return h.refresh(bot, blackList, query)
	case model.CallbackTypePassThrough:
		blackList := &model.BlackList{}
		err := db.Instance().Model(blackList).
			Where("group_id = ?", query.Message.Chat.ID).
			Where("captcha_msg_id = ?", query.Message.MessageID).
			First()
		if err != nil {
			_, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, err.Error()))
			return err
		}
		return h.passThrough(bot, blackList, query)
	case model.CallbackTypeDonateWX, model.CallbackTypeDonateAlipay:
		return extra.UpdateMsgPhoto(bot, query.Message.Chat.ID, query.Message.MessageID, "", "",
			model.DonatesKeyboard(query.Data), model.Donates[query.Data].ID)
	}
	return nil
}

func (h Handler) refresh(bot *tgbotapi.BotAPI, blackList *model.BlackList, query tgbotapi.CallbackQuery) error {
	if query.From.ID != blackList.UserId {
		return h.answerCallbackQuery(bot, query, "无权限")
	}
	idiom, err := h.GetRandomIdiom()
	if err != nil {
		h.answerCallbackQuery(bot, query, "刷新失败")
		return err
	}
	blackList.IdiomId = idiom.Id
	if _, err := db.Instance().Model(blackList).
		Column("idiom_id").
		WherePK().Update(); err != nil {
		h.answerCallbackQuery(bot, query, "刷新失败")
		return err
	}
	if err := extra.UpdateMsgPhoto(
		bot, query.Message.Chat.ID, query.Message.MessageID,
		fmt.Sprintf(model.EnterRoomMsg, blackList.UserLink, query.Message.Chat.Title,
			time.Until(blackList.ExpireAt)/time.Second),
		tgbotapi.ModeMarkdown, model.InlineKeyboard, tgbotapi.FileBytes{
			Name:  strconv.Itoa(query.From.ID),
			Bytes: idiom.CaptchaImg,
		},
	); err != nil {
		h.answerCallbackQuery(bot, query, "刷新失败")
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
	_, err := db.Instance().Model(blackList).
		Where("group_id = ?group_id").
		Where("user_id = ?user_id").
		Delete()
	if err != nil {
		return err
	}
	bot.DeleteMessage(tgbotapi.NewDeleteMessage(blackList.GroupId, blackList.CaptchaMsgId))
	passMsgConfig := tgbotapi.NewMessage(blackList.GroupId, fmt.Sprintf(
		"%s 恭喜，你已验证通过", blackList.UserLink))
	passMsgConfig.ParseMode = tgbotapi.ModeMarkdown
	passMsg, err := bot.Send(passMsgConfig)
	if err != nil {
		return err
	}
	scheduler.UpdateMsgExpireTaskDone(db.Instance(), blackList.Id)
	return scheduler.AddDelMsgTask(db.Instance(), blackList.GroupId, passMsg.MessageID)
}

func (h Handler) GetRandomIdiom() (*model.Idiom, error) {
	idiom := &model.Idiom{}
	rand.Seed(time.Now().UnixNano())
	randOffset := rand.Intn(h.IdiomCount)
	if err := db.Instance().Model(idiom).
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
