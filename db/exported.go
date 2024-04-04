package db

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/jqs7/zwei/env"
)

func WithEnv(cfg env.Specification) OptionFunc {
	return func(op *pg.Options) {
		op.Addr = fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
		op.User = cfg.User
		op.Password = cfg.Password
		op.Database = cfg.Database
		op.OnConnect = func(db *pg.DB) error {
			log.Println("database is connected")
			return nil
		}
	}
}
