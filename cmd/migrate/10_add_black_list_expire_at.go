package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists ADD expire_at timestamptz NULL;")
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists DROP COLUMN expire_at;")
		return err
	})
}
