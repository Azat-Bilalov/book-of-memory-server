package app

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/routes"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	e := echo.New()

	e.Use(middleware.Recover())

	userRepository := repository.NewUserRepository(config.DB)
	veteranRepository := repository.NewVeteranRepository(config.DB)
	docBindingRepository := repository.NewDocBindingRepository(config.DB)

	bindingRepository := repository.NewBindingRepository(config.DB)
	bindingUsecase := usecase.NewBindingUsecase(bindingRepository, userRepository, veteranRepository)
	bindingHandler := handler.NewBindingHandler(bindingUsecase)
	routes.InitBindingRoutes(e, bindingHandler)

	documentRepository := repository.NewDocumentRepository(config.DB)
	documentUsecase := usecase.NewDocumentUsecase(documentRepository, bindingRepository, docBindingRepository)
	documentHandler := handler.NewDocumentHandler(documentUsecase)
	routes.InitDocumentRoutes(e, documentHandler)

	e.Logger.Fatal(e.Start(":8080"))
	log.Println("Server down")
}
