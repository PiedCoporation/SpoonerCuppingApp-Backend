package main

import (
	"backend/config"
	"backend/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
