package main

import (
	"os"

	"github.com/go-pg/pg/orm"

	"github.com/json-iterator/go"

	"github.com/go-pg/migrations"
)

type Idiom struct {
	ID           int64
	Derivation   string
	Example      string
	Explanation  string
	Pinyin       string
	Word         string
	Abbreviation string
}

func init() {
	migrations.Register(func(db migrations.DB) error {
		if err := db.Model(new(Idiom)).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			return err
		}
		f, err := os.Open("./idiom.json")
		if err != nil {
			return err
		}
		var idioms []Idiom
		if err := jsoniter.NewDecoder(f).Decode(&idioms); err != nil {
			return err
		}
		for i := range idioms {
			if err := db.Insert(&idioms[i]); err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		return db.Model(new(Idiom)).DropTable(&orm.DropTableOptions{
			IfExists: true,
		})
	})
}
