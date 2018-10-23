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

func Instance() *pg.DB {
	once.Do(func() {
		pgDB = pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", env.Spec.Address, env.Spec.Port),
			User:     env.Spec.User,
			Database: env.Spec.Database,
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
