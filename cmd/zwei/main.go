package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	runScheduler(ctx, wg, bot)
	runBot(ctx, wg, bot)
	gracefulListener(bot, cancel)
	wg.Wait()
}

func runScheduler(ctx context.Context, wg *sync.WaitGroup, bot *tg.Bot) {
	go func() {
		err := scheduler.New(db.Instance(), bot.BotAPI).Run(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()
}

func runBot(ctx context.Context, wg *sync.WaitGroup, bot *tg.Bot) {
	go func() {
		bot.Run(ctx)
		wg.Done()
	}()
}

func gracefulListener(bot *tg.Bot, cancel context.CancelFunc) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		bot.StopReceivingUpdates()
		cancel()
	}()
}
