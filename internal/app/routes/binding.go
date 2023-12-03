package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitBindingRoutes(e *echo.Echo, bindingHandler *handler.BindingHandler, m *Middlewares) {
	e.GET("/bindings", bindingHandler.FindBindings, m.WithAuth)
	e.GET("/bindings/:uuid", bindingHandler.FindBindingByUUID, m.WithAuth)
	e.PUT("/bindings/:uuid", bindingHandler.UpdateBindingByUUID, m.WithAuth, m.WithUser)
	e.PUT("/bindings/:uuid/submit", bindingHandler.SubmitBindingByUUID, m.WithAuth, m.WithUser)
	e.PUT("/bindings/:uuid/accept-reject", bindingHandler.AcceptRejectBindingByUUID, m.WithAuth, m.WithAdmin)
	e.DELETE("/bindings/:uuid", bindingHandler.DeleteBindingByUUID, m.WithAuth, m.WithUser)
}
