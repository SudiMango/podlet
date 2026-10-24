package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/SudiMango/podlet/internal/agent"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	app, err := agent.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
