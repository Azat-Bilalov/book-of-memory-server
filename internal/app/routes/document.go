package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitDocumentRoutes(e *echo.Echo, documentHandler *handler.DocumentHandler, m *Middlewares) {
	e.GET("/documents", documentHandler.FindActiveDocuments, m.WithOptionalAuth)
	e.GET("/documents/:uuid", documentHandler.FindActiveDocumentByUUID)
	e.POST("/documents", documentHandler.CreateDocument)
	e.PUT("/documents/:uuid", documentHandler.UpdateDocumentByUUID)
	e.DELETE("/documents/:uuid", documentHandler.DeleteDocumentByUUID)
	e.POST("/bindings/documents/:uuid", documentHandler.AddDocumentToBindingByUUID)
	e.DELETE("/bindings/documents/:uuid", documentHandler.RemoveDocumentFromBindingByUUID)
}
