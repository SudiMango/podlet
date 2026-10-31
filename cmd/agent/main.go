package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/SudiMango/podlet/internal/agent"
)

func main() {
	initializeOnly := flag.Bool("initialize-only", false, "initialize the agent and exit")
	flag.Parse()

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

	if *initializeOnly {
		if err := app.Initialize(ctx, true); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Podlet initialization completed!")
		return
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
