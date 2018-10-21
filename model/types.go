package model

import (
	"time"
)

type BlackList struct {
	Id           int64
	GroupId      int64
	UserId       int
	IdiomId      int64
	Idiom        Idiom
	CaptchaMsgId int
	UserLink     string
	DeletedAt    *time.Time `pg:",soft_delete"`
}

type Idiom struct {
	Id           int64
	Derivation   string
	Example      string
	Explanation  string
	Pinyin       string
	Word         string
	Abbreviation string
	CaptchaImg   []byte `sql:"-"`
}

const (
	TaskTypeDeleteMsg = "DeleteMsg"
)

const (
	TaskStatusPlan = iota
	TaskStatusDoing
	TaskStatusDone
)

type Task struct {
	Id     int64
	Type   string
	Status int64 `sql:",notnull"`
	RunAt  time.Time
	ChatID int64
	MsgID  int
}
