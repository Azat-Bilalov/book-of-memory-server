package handler

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
)

type BindingHandler struct {
	bindingUsecase *usecase.BindingUsecase
}

func NewBindingHandler(bindingUsecase *usecase.BindingUsecase) *BindingHandler {
	return &BindingHandler{bindingUsecase}
}

func (h *BindingHandler) FindBindings(c echo.Context) error {
	var (
		userID   = c.Request().Header.Get("x-user-id")
		status   = c.QueryParam("status")
		dateFrom = c.QueryParam("date_from")
		dateTo   = c.QueryParam("date_to")
	)
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	bindings, err := h.bindingUsecase.FindBindingsByUserID(userID, status, dateFrom, dateTo)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, bindings)
}

func (h *BindingHandler) FindBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, "заявка не найдена")
	}
	if binding.UserID != userID && binding.ModeratorID != userID {
		return c.NoContent(403)
	}
	return c.JSONPretty(200, binding, " ")
}

func (h *BindingHandler) UpdateBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, "заявка не найдена")
	}
	if binding.UserID != userID && binding.ModeratorID != userID {
		return c.NoContent(403)
	}
	if binding.Status != ds.BINDING_STATUS_ENTERED && binding.Status != ds.BINDING_STATUS_IN_PROGRESS {
		return c.JSON(400, "статус не является entered или in_progress")
	}
	req := ds.BindingUpdateRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, err.Error())
	}
	binding, err = h.bindingUsecase.UpdateBindingByUUID(uuid, req)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, binding)
}

func (h *BindingHandler) SubmitBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, "заявка не найдена")
	}
	if binding.UserID != userID {
		return c.NoContent(403)
	}
	if binding.Status != ds.BINDING_STATUS_ENTERED {
		return c.JSON(400, "статус не является entered")
	}
	if binding.VeteranID == nil {
		return c.JSON(400, "ветеран не установлен")
	}
	binding, err = h.bindingUsecase.SubmitBindingByUUID(uuid)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, binding)
}

func (h *BindingHandler) AcceptRejectBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	moderatorID := c.Request().Header.Get("x-user-id")
	if moderatorID == "" || !IsValidUUID(moderatorID) {
		return c.JSON(400, "идентификатор модератора пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, "заявка не найдена")
	}
	if binding.ModeratorID != moderatorID {
		return c.NoContent(403)
	}
	if binding.Status != ds.BINDING_STATUS_IN_PROGRESS {
		return c.JSON(400, "статус не является in_progress")
	}
	req := ds.BindingStatusUpdateRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, err.Error())
	}
	switch req.Status {
	case ds.BINDING_STATUS_COMPLETED:
		binding, err = h.bindingUsecase.AcceptBindingByUUID(uuid)
	case ds.BINDING_STATUS_CANCELED:
		binding, err = h.bindingUsecase.RejectBindingByUUID(uuid)
	default:
		return c.JSON(400, "статус недействителен")
	}
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, binding)
}

func (h *BindingHandler) DeleteBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, err.Error())
	}
	if binding.UserID != userID {
		return c.NoContent(403)
	}
	if binding.Status != ds.BINDING_STATUS_ENTERED {
		return c.JSON(400, "статус не является entered")
	}
	err = h.bindingUsecase.DeleteBindingByUUID(uuid)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}
