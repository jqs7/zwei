package scheduler

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/model"
)

type Scheduler struct {
	*pg.DB
	*tgbotapi.BotAPI
}

func New(db *pg.DB, bot *tgbotapi.BotAPI) *Scheduler {
	db.Model(&model.Task{}).
		Where("status = ?", model.TaskStatusDoing).
		Set("status = ?", model.TaskStatusPlan).
		Update()
	return &Scheduler{
		DB:     db,
		BotAPI: bot,
	}
}

func (s Scheduler) Run() error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		var tasks []model.Task
		err := s.Model(&model.Task{}).
			Where("status = ?", model.TaskStatusPlan).
			Where("run_at < ?", time.Now()).
			Limit(10).
			Select(&tasks)
		if err != nil {
			return err
		}
		for _, task := range tasks {
			s.processTask(task)
		}
	}
	return nil
}

func (s Scheduler) processTask(task model.Task) error {
	_, err := s.Model(&task).WherePK().
		Set("status = ?", model.TaskStatusDoing).
		Update()
	if err != nil {
		return err
	}
	switch task.Type {
	case model.TaskTypeDeleteMsg:
		s.DeleteMessage(tgbotapi.NewDeleteMessage(task.ChatID, task.MsgID))
	}
	_, err = s.Model(&task).WherePK().
		Set("status = ?", model.TaskStatusDone).
		Update()
	return err
}

func AddDelMsgTask(db *pg.DB, chatID int64, msgID int) error {
	return db.Insert(&model.Task{
		Type:   model.TaskTypeDeleteMsg,
		Status: model.TaskStatusPlan,
		RunAt:  time.Now().Add(time.Second * 10),
		ChatID: chatID,
		MsgID:  msgID,
	})
}
