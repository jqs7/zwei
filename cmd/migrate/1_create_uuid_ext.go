package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP EXTENSION "uuid-ossp";`)
		return err
	})
}
