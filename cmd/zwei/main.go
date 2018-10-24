package main

import (
	"log"

	"github.com/jqs7/zwei/biz"
	"github.com/jqs7/zwei/bot/tg"
	"github.com/jqs7/zwei/db"
	"github.com/jqs7/zwei/env"
	"github.com/jqs7/zwei/scheduler"
)

func main() {
	env.Init()

	var botOpts []tg.ConfigFunc
	if env.Spec.Debug {
		botOpts = append(botOpts, tg.Debug())
	}
	handler := biz.NewHandler(180)
	bot := tg.NewBot(env.Spec.Token, handler, botOpts...)
	go func() {
		err := scheduler.New(db.Instance(), bot.BotAPI).Run()
		log.Panic(err)
	}()
	bot.Run()
}
