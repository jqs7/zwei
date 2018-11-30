package main

import (
	"encoding/json"
	"os"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"
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
	_ = migrations.Register(func(db migrations.DB) error {
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
		if err := json.NewDecoder(f).Decode(&idioms); err != nil {
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
