package app

import (
	"log"

	_ "github.com/Azat-Bilalov/book-of-memory-server/docs"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/repository"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/routes"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	withAuth := a.WithAuth
	WithOptionalAuth := a.WithOptionalAuth
	withAdmin := func(next echo.HandlerFunc) echo.HandlerFunc {
		return a.WithAuth(a.WithRole("moderator", next))
	}
	withUser := func(next echo.HandlerFunc) echo.HandlerFunc {
		return a.WithAuth(a.WithRole("user", next))
	}

	m := &routes.Middlewares{
		WithAuth:         withAuth,
		WithOptionalAuth: WithOptionalAuth,
		WithAdmin:        withAdmin,
		WithUser:         withUser,
	}

	e := echo.New()

	// e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	userRepository := repository.NewUserRepository(a.repository.DB)
	veteranRepository := repository.NewVeteranRepository(a.repository.DB)
	docBindingRepository := repository.NewDocBindingRepository(a.repository.DB)

	bindingRepository := repository.NewBindingRepository(a.repository.DB)
	bindingUsecase := usecase.NewBindingUsecase(bindingRepository, userRepository, veteranRepository)
	bindingHandler := handler.NewBindingHandler(bindingUsecase)
	routes.InitBindingRoutes(e, bindingHandler, m)

	documentRepository := repository.NewDocumentRepository(a.repository.DB)
	documentUsecase := usecase.NewDocumentUsecase(documentRepository, bindingRepository, docBindingRepository, userRepository)
	documentHandler := handler.NewDocumentHandler(documentUsecase)
	routes.InitDocumentRoutes(e, documentHandler, m)

	routes.InitFileRoutes(e, handler.NewFileHandler())

	authUsecase := usecase.NewAuthUsecase(userRepository, a.redis)
	authHandler := handler.NewAuthHandler(authUsecase)
	routes.InitRegisterRoutes(e, authHandler, m)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8080"))
	log.Println("Server down")
}
