package main

import (
	"time"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"
)

type Task struct {
	Id     int64
	Type   string
	Status int64
	RunAt  time.Time
	ChatID int64
	MsgID  int
}

func init() {
	migrations.Register(func(db migrations.DB) error {
		return db.Model(&Task{}).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
	}, func(db migrations.DB) error {
		return db.Model(&Task{}).DropTable(&orm.DropTableOptions{
			IfExists: true,
		})
	})
}
