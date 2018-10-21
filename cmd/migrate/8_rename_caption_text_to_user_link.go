package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE public.black_lists RENAME COLUMN caption_text TO user_link;")
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE public.black_lists RENAME COLUMN user_linkcaption_text TO caption_text;")
		return err
	})
}
