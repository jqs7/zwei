package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/jqs7/zwei/db"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/bot/extra"
	"github.com/jqs7/zwei/model"
)

type Scheduler struct {
	*db.DB
	*tgbotapi.BotAPI
}

func New(db *db.DB, bot *tgbotapi.BotAPI) *Scheduler {
	db.PgDB.Model(&model.Task{}).
		Where("status = ?", model.TaskStatusDoing).
		Set("status = ?", model.TaskStatusPlan).
		Update()
	return &Scheduler{
		DB:     db,
		BotAPI: bot,
	}
}

func (s Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var tasks []model.Task
			err := s.PgDB.Model(&model.Task{}).
				Where("status = ?", model.TaskStatusPlan).
				Where("run_at <= ?", time.Now()).
				Limit(10).
				Select(&tasks)
			if err != nil {
				return err
			}
			for _, task := range tasks {
				s.processTask(task)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s Scheduler) processTask(task model.Task) error {
	_, err := s.PgDB.Model(&task).WherePK().
		Set("status = ?", model.TaskStatusDoing).
		Update()
	if err != nil {
		return err
	}
	switch task.Type {
	case model.TaskTypeDeleteMsg:
		s.DeleteMessage(tgbotapi.NewDeleteMessage(task.ChatID, task.MsgID))
		_, err = s.PgDB.Model(&task).WherePK().
			Set("status = ?", model.TaskStatusDone).
			Update()
		return err
	case model.TaskTypeUpdateMsgExpire:
		return s.updateMsgExpire(task)
	}
	return nil
}

func (s Scheduler) updateMsgExpire(task model.Task) error {
	blackList := &model.BlackList{Id: task.BlackListId}
	err := s.PgDB.Model(blackList).WherePK().First()
	if err != nil {
		return s.taskDone(&task)
	}
	timeSub := blackList.ExpireAt.Sub(time.Now()) / time.Second
	if timeSub <= 0 {
		extra.KickAndDelCaptcha(s.BotAPI, *blackList, time.Now().Add(time.Minute).Unix())
		return s.taskDone(&task)
	}
	if err := s.updateMsg(blackList, timeSub); err != nil {
		return s.taskDelay(&task, model.DefaultRefreshDuration)
	}
	return s.taskDelay(&task, model.DefaultRefreshDuration)
}

func (s Scheduler) taskDone(task *model.Task) error {
	_, err := s.PgDB.Model(task).
		WherePK().
		Set("status = ?", model.TaskStatusDone).
		Update()
	return err
}

func (s Scheduler) taskDelay(task *model.Task, dur time.Duration) error {
	_, err := s.PgDB.Model(task).
		WherePK().
		Set("status = ?", model.TaskStatusPlan).
		Set("run_at = ?", time.Now().Add(dur)).
		Update()
	return err
}

func (s Scheduler) updateMsg(blackList *model.BlackList, timeSub time.Duration) error {
	chat, err := s.GetChat(tgbotapi.ChatConfig{ChatID: blackList.GroupId})
	if err != nil {
		return err
	}
	caption := fmt.Sprintf(model.EnterRoomMsg, blackList.UserLink, chat.Title, timeSub)
	editor := tgbotapi.NewEditMessageCaption(blackList.GroupId, blackList.CaptchaMsgId, caption)
	editor.ReplyMarkup = &model.InlineKeyboard
	editor.ParseMode = tgbotapi.ModeHTML
	_, err = s.Send(editor)
	return err
}

func AddDelMsgTask(db *db.DB, chatID int64, msgID int) error {
	return db.PgDB.Insert(&model.Task{
		Type:   model.TaskTypeDeleteMsg,
		Status: model.TaskStatusPlan,
		RunAt:  time.Now().Add(time.Second * 10),
		ChatID: chatID,
		MsgID:  msgID,
	})
}

func AddUpdateMsgExpireTask(db *db.DB, blackListID, chatID int64, msgID int) error {
	return db.PgDB.Insert(&model.Task{
		Type:        model.TaskTypeUpdateMsgExpire,
		Status:      model.TaskStatusPlan,
		RunAt:       time.Now().Add(model.DefaultRefreshDuration),
		ChatID:      chatID,
		MsgID:       msgID,
		BlackListId: blackListID,
	})
}

func UpdateMsgExpireTaskDone(db *db.DB, blackListID int64) error {
	_, err := db.PgDB.Model(&model.Task{}).
		Where("type = ?", model.TaskTypeUpdateMsgExpire).
		Where("black_list_id = ?", blackListID).
		Set("status = ?", model.TaskStatusDone).
		Update()
	return err
}
