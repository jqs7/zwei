// +build wireinject

package main

import (
	"context"
	"sync"

	"github.com/google/wire"
	"github.com/jqs7/zwei/biz"
	"github.com/jqs7/zwei/db"
	"github.com/jqs7/zwei/env"
)

func Run(ctx context.Context, cancel context.CancelFunc) *sync.WaitGroup {
	wire.Build(
		env.Init,
		db.Instance,
		biz.NewHandler,
		ProvideBot,
		Runner,
	)
	return &sync.WaitGroup{}
}
