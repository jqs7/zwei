package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-pg/migrations"
	"github.com/jqs7/zwei/model"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		type BlackList struct {
			Id           int64
			GroupId      int64
			UserId       int
			IdiomId      int64
			Idiom        Idiom
			CaptchaMsgId int
			UserLink     string
			ExpireAt     time.Time
			DeletedAt    *time.Time `pg:",soft_delete"`
		}

		var bls []BlackList
		if err := db.Model(&model.BlackList{}).Select(&bls); err != nil {
			return err
		}
		for _, bl := range bls {
			matches := regexp.MustCompile(`\[(.+)]\(tg://user\?id=(\d+)\)`).FindStringSubmatch(bl.UserLink)
			if len(matches) < 3 {
				continue
			}
			db.Model(&bl).WherePK().Set(
				"user_link = ?",
				fmt.Sprintf(`<a href="tg://user?id=%s">%s</a>`, matches[2], matches[1]),
			).Update()
		}
		return nil
	}, func(db migrations.DB) error {
		return nil
	})
}
