package main

import (
	"context"
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app"
)

// @title Book of Memory API
// @version 1.0
// @description This is a sample server Book of Memory API server.

// @contact.name API Support
// @contact.url https://t.me/azat_bil
// @contact.email az@bilalov@mail.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1:8080
// @schemes http
// @BasePath /

// @SecurityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

func main() {
	log.Println("Application start!")

	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	application.StartServer()

	log.Println("Application terminated!")
}
