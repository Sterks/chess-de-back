package main

import (
	"chess-backend/internal/app"
	"log"
)

const configsDir = "../configs"

func main() {

	if err := app.Run(configsDir); err != nil {
		log.Fatalln(err)
	}

}
