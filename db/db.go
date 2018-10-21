package db

import (
	"log"
	"sync"
	"time"

	"github.com/go-pg/pg"
)

var (
	pgDB *pg.DB
	once sync.Once
)

func Instance() *pg.DB {
	once.Do(func() {
		pgDB = pg.Connect(&pg.Options{
			User:     "jqs7",
			Database: "zwei",
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
