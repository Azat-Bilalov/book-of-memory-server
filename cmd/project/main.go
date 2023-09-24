package main

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app"
)

func main() {
	log.Println("Application start!")

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	application.StartServer()

	log.Println("Application terminated!")
}
