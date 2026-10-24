package main

import (
	"github.com/SudiMango/podlet/internal/cli"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cli.RootCommand{}.Execute()
}
