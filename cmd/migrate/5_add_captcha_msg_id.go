package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists ADD captcha_msg_id int8 NULL;")
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists DROP COLUMN captcha_msg_id;")
		return err
	})
}
