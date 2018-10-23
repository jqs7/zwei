package main

import (
	"log"
	"os"

	"github.com/hanguofeng/gocaptcha"
	"github.com/jqs7/zwei/biz"
	"github.com/jqs7/zwei/bot/tg"
	"github.com/jqs7/zwei/db"
	"github.com/jqs7/zwei/env"
	"github.com/jqs7/zwei/model"
	"github.com/jqs7/zwei/scheduler"
)

func main() {

	err := env.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	filterConfig := new(gocaptcha.FilterConfig)
	filterConfig.Init()
	filterConfig.Filters = []string{
		gocaptcha.IMAGE_FILTER_NOISE_LINE,
		gocaptcha.IMAGE_FILTER_NOISE_POINT,
		gocaptcha.IMAGE_FILTER_STRIKE,
	}
	for _, v := range filterConfig.Filters {
		filterConfigGroup := new(gocaptcha.FilterConfigGroup)
		filterConfigGroup.Init()
		filterConfigGroup.SetItem("Num", "180")
		filterConfig.SetGroup(v, filterConfigGroup)
	}
	idiomCount, err := db.Instance().Model(new(model.Idiom)).Count()
	if err != nil {
		log.Fatalln(err)
	}

	pwd, _ := os.Getwd()
	fontPath := pwd + "/fonts/"
	bot := tg.NewBot(
		env.Spec.Token,
		biz.Handler{
			ImageConfig: &gocaptcha.ImageConfig{
				Width:    320,
				Height:   100,
				FontSize: 80,
				FontFiles: []string{
					fontPath + "STFANGSO.ttf",
					fontPath + "STHEITI.ttf",
					fontPath + "STXIHEI.ttf",
				},
			},
			IdiomCount:         idiomCount,
			ImageFilterManager: gocaptcha.CreateImageFilterManagerByConfig(filterConfig),
		},
		tg.Debug(),
	)
	go func() {
		err := scheduler.New(db.Instance(), bot.BotAPI).Run()
		log.Panic(err)
	}()
	bot.Run()
}
