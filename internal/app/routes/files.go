package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitFileRoutes(e *echo.Echo, fileHandler *handler.FileHandler) {
	e.GET("/files/:bucket/:filename", fileHandler.FindFile)
}
