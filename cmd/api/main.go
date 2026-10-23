package main

import (
	"log"
	"log/slog"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/handler"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/server"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/database"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	slog.Info("starting backend...")

	fx.New(
		database.Module,
		auth.Module,
		repository.Module,
		service.Module,
		handler.Module,
		server.Module,
	).Run()
}
