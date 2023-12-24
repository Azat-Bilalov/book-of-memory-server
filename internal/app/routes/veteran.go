package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitVeteranRoutes(e *echo.Echo, veteranHadler *handler.VeteranHandler, m *Middlewares) {
	e.GET("/veterans", veteranHadler.FindVeterans, m.WithOptionalAuth)
}
