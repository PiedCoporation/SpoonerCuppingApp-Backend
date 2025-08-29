package main

import (
	"backend/config"
	"fmt"
	"log"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// test load config
	fmt.Println(config)
}
