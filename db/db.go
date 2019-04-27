package db

import (
	"log"
	"sync"
	"time"

	"github.com/go-pg/pg"
)

var (
	db   *DB
	once sync.Once
)

type DB struct {
	PgDB *pg.DB
}

type OptionFunc func(*pg.Options)

func Instance(opts ...OptionFunc) *DB {
	once.Do(func() {
		pgOpt := new(pg.Options)
		for _, opt := range opts {
			opt(pgOpt)
		}
		pgDB := pg.Connect(pgOpt)
		pgDB.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
		db = &DB{PgDB: pgDB}
	})
	return db
}
