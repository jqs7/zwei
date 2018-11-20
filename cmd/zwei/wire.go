// +build wireinject

package main

import (
	"context"
	"sync"

	"github.com/google/go-cloud/wire"
	"github.com/jqs7/zwei/biz"
	"github.com/jqs7/zwei/env"
)

func Run(ctx context.Context, cancel context.CancelFunc) *sync.WaitGroup {
	wire.Build(
		env.Init,
		ProvideDB,
		biz.NewHandler,
		ProvideBot,
		Runner,
	)
	return &sync.WaitGroup{}
}
