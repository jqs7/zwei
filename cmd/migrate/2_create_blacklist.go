package main

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"
)

type BlackList struct {
	Id      int64
	GroupId int64
	UserId  int
}

func init() {
	migrations.Register(func(db migrations.DB) error {
		return db.Model(&BlackList{}).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
	}, func(db migrations.DB) error {
		return db.Model(&BlackList{}).DropTable(&orm.DropTableOptions{
			IfExists: true,
		})
	})
}
