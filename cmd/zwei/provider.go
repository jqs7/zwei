package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-pg/pg"
	"github.com/jqs7/zwei/biz"
	"github.com/jqs7/zwei/bot/tg"
	"github.com/jqs7/zwei/env"
	"github.com/jqs7/zwei/scheduler"
)

func ProvideBot(cfg env.Specification, handler biz.Handler) *tg.Bot {
	var botOpts []tg.ConfigFunc
	if cfg.Debug {
		botOpts = append(botOpts, tg.Debug())
	}
	return tg.NewBot(env.Spec.Token, handler, botOpts...)
}

func Runner(ctx context.Context, cancel context.CancelFunc, db *pg.DB, bot *tg.Bot) *sync.WaitGroup {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		bot.StopReceivingUpdates()
		cancel()
	}()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		bot.Run(ctx)
		wg.Done()
	}()
	go func() {
		err := scheduler.New(db, bot.BotAPI).Run(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()
	return wg
}
