package main

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
