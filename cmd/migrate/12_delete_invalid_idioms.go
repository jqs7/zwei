package main

import (
	"log"

	"github.com/go-pg/migrations"
	"github.com/jqs7/zwei/model"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		rst, err := db.Model(&model.Idiom{}).
			Where("char_length(word) <> ?", 4).
			Delete()
		log.Printf("%d rows deleted", rst.RowsAffected())
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("")
		return err
	})
}
