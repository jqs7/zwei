package main

import (
	"context"
)

func main() {
	Run(context.WithCancel(context.Background())).Wait()
}
