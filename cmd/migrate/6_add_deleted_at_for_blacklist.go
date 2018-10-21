package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists ADD deleted_at timestamp NULL;")
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec("ALTER TABLE black_lists DROP COLUMN deleted_at;")
		return err
	})
}
