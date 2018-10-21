package tg

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jqs7/zwei/biz"
)

type RunningConfig struct {
	debug bool
}

type ConfigFunc func(*RunningConfig)

func Debug() ConfigFunc {
	return func(config *RunningConfig) {
		config.debug = true
	}
}

func NewBot(botAPI string, handler biz.IBiz, fs ...ConfigFunc) *Bot {
	cfg := &RunningConfig{
		debug: false,
	}
	for _, f := range fs {
		f(cfg)
	}
	bot, err := tgbotapi.NewBotAPI(botAPI)
	if err != nil {
		log.Fatalln(err)
	}
	if cfg.debug {
		bot.Debug = true
	}

	myInfo, err := bot.GetMe()
	if err != nil {
		log.Fatalln(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln(err)
	}

	return &Bot{
		BotAPI:  bot,
		myInfo:  myInfo,
		IBiz:    handler,
		updates: updates,
	}
}

func (b Bot) Run() {
	for update := range b.updates {
		err := b.processUpdate(update)
		if err != nil {
			log.Printf("got err: %+v\n", err)
		}
	}
}
