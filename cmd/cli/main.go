package main

import (
	"log"

	"github.com/SudiMango/podlet/internal/cli"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	cli.Execute()
}
