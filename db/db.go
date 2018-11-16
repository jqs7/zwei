package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-pg/pg"
	"github.com/jqs7/zwei/env"
)

var (
	pgDB *pg.DB
	once sync.Once
)

func Instance(cfg env.Specification) *pg.DB {
	once.Do(func() {
		pgDB = pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Address, cfg.Port),
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
			OnConnect: func(db *pg.DB) error {
				log.Println("database is connected")
				return nil
			},
		})
		pgDB.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
	})
	return pgDB
}
