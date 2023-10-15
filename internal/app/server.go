package app

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/routes"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	e := echo.New()

	documentRepository := repository.NewDocumentRepository(config.DB)
	documentUsecase := usecase.NewDocumentUsecase(documentRepository)
	documentHandler := handler.NewDocumentHandler(documentUsecase)
	routes.InitDocumentRoutes(e, documentHandler)

	e.Logger.Fatal(e.Start(":8080"))
	log.Println("Server down")
}
