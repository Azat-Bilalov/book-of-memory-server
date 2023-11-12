package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitBindingRoutes(e *echo.Echo, bindingHandler *handler.BindingHandler) {
	e.GET("/bindings", bindingHandler.FindBindings)
	e.GET("/bindings/:uuid", bindingHandler.FindBindingByUUID)
	e.PUT("/bindings/:uuid", bindingHandler.UpdateBindingByUUID)
	e.PUT("/bindings/:uuid/submit", bindingHandler.SubmitBindingByUUID)
	e.PUT("/bindings/:uuid/accept-reject", bindingHandler.AcceptRejectBindingByUUID)
	e.DELETE("/bindings/:uuid", bindingHandler.DeleteBindingByUUID)
}
