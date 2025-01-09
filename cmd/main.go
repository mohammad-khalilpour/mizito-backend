package main

import (
	"fmt"
	config "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"mizito/internal/env"
	"mizito/internal/router"
	"os"
)

func main() {

	var cfg env.Config

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	if err := config.Parse(&cfg); err != nil {
		panic(err)
	}

	zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	r := router.InitApp(&cfg)
	fmt.Println("initializing routes ...")
	r.Init(&cfg)
	fmt.Println("initializing completed...")

	fmt.Println("running server...")
	r.Run()

}
