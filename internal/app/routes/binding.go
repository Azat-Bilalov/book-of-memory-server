package routes

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/handler"
	"github.com/labstack/echo/v4"
)

func InitBindingRoutes(e *echo.Echo, bindingHandler *handler.BindingHandler) {
	e.GET("/bindings", bindingHandler.FindBindings)
	e.GET("/bindings/:uuid", bindingHandler.FindBindingByUUID)
	e.PUT("/bindings/:uuid", bindingHandler.UpdateBindingByUUID)        // пользователь/модератор обновляют информацию в заявке
	e.PUT("/bindings/:uuid/submit", bindingHandler.SubmitBindingByUUID) // пользователь отпраляет заявку на рассмотрение модератору
	e.PUT("/bindings/:uuid/accept", bindingHandler.AcceptBindingByUUID) // модератор принимает заявку
	e.PUT("/bindings/:uuid/reject", bindingHandler.RejectBindingByUUID) // модератор отклоняет заявку
	e.DELETE("/bindings/:uuid", bindingHandler.DeleteBindingByUUID)     // пользователь/модератор удаляют заявку
}
