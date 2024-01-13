package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitVeteranRoutes(e *echo.Echo, veteranHadler *handler.VeteranHandler, m *Middlewares) {
	e.GET("/veterans", veteranHadler.FindVeterans, m.WithOptionalAuth)
	e.POST("/veterans", veteranHadler.CreateVeteran, m.WithAuth)
	e.GET("/veterans/:uuid", veteranHadler.FindVeteranByUUID, m.WithOptionalAuth)
	e.PUT("/veterans/:uuid", veteranHadler.UpdateVeteranByUUID, m.WithAuth, m.WithAdmin)
	e.DELETE("/veterans/:uuid", veteranHadler.DeleteVeteranByUUID, m.WithAuth, m.WithAdmin)
}
