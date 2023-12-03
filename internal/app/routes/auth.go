package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitRegisterRoutes(e *echo.Echo, authHadler *handler.AuthHandler, m *Middlewares) {
	e.POST("/login", authHadler.Login)
	e.POST("/register", authHadler.Register)
	e.POST("/logout", authHadler.Logout, m.WithAuth)
}
