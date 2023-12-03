package handler

import (
	"strings"

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

// Binding godoc
// @Summary      Find bindings
// @Description  Find bindings
// @Tags         Binding
// @Produce      json
// @Param        status   query string false "Status"
// @Param        date_from query string false "Date from" Format(date)
// @Param        date_to   query string false "Date to" Format(date)
// @Success      200  {object}  []ds.Binding
// @Failure      400  {string}  string "некорректный запрос"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Security     JwtAuth
// @Router       /bindings [get]
func (h *BindingHandler) FindBindings(c echo.Context) error {
	var (
		userID   = c.Get("user_id").(string)
		status   = strings.ToLower(c.QueryParam("status"))
		dateFrom = strings.ToLower(c.QueryParam("date_from"))
		dateTo   = strings.ToLower(c.QueryParam("date_to"))
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

// Binding godoc
// @Summary      Find binding by UUID
// @Description  Find binding by UUID
// @Tags         Binding
// @Produce      json
// @Param        uuid path string true "UUID"
// @Success      200  {object}  ds.Binding
// @Failure      400  {string}  string "недействительный идентификатор"
// @Failure      400  {string}  string "идентификатор пользователя пустой или недействительный"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Failure      404  {string}  string "заявка не найдена"
// @Security     JwtAuth
// @Router       /bindings/{uuid} [get]
func (h *BindingHandler) FindBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Get("user_id").(string)
	userRole := c.Get("user_role").(string)
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или недействительный")
	}
	binding, err := h.bindingUsecase.FindBindingByUUID(uuid)
	if err != nil {
		return c.JSON(404, "заявка не найдена")
	}
	if binding.UserID == userID || userRole == "moderator" {
		return c.NoContent(403)
	}
	return c.JSONPretty(200, binding, " ")
}

// Binding godoc
// @Summary      Update binding by UUID
// @Description  Update binding by UUID
// @Tags         Binding
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID"
// @Param        body body ds.BindingUpdateRequest true "Body"
// @Success      200  {object}  ds.Binding
// @Failure      400  {string}  string "недействительный идентификатор"
// @Failure      400  {string}  string "идентификатор пользователя пустой или недействительный"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Failure      404  {string}  string "заявка не найдена"
// @Failure      400  {string}  string "статус не является entered или in_progress"
// @Security     JwtAuth
// @Router       /bindings/{uuid} [put]
func (h *BindingHandler) UpdateBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Get("user_id").(string)
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

// Binding godoc
// @Summary      Submit binding by UUID
// @Description  Submit binding by UUID
// @Tags         Binding
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID"
// @Success      200  {object}  ds.Binding
// @Failure      400  {string}  string "недействительный идентификатор"
// @Failure      400  {string}  string "идентификатор пользователя пустой или недействительный"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Failure      404  {string}  string "заявка не найдена"
// @Failure      400  {string}  string "статус не является entered"
// @Failure      400  {string}  string "ветеран не установлен"
// @Security     JwtAuth
// @Router       /bindings/{uuid}/submit [put]
func (h *BindingHandler) SubmitBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "недействительный идентификатор")
	}
	userID := c.Get("user_id").(string)
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

// Binding godoc
// @Summary      Accept or reject binding by UUID
// @Description  Accept or reject binding by UUID
// @Tags         Binding
// @Accept       json
// @Produce      json
// @Param        uuid path string true "UUID"
// @Param        body body ds.BindingStatusUpdateRequest true "Body"
// @Success      200  {object}  ds.Binding
// @Failure      400  {string}  string "недействительный идентификатор"
// @Failure      400  {string}  string "идентификатор модератора пустой или недействительный"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Failure      404  {string}  string "заявка не найдена"
// @Failure      400  {string}  string "статус не является in_progress"
// @Failure      400  {string}  string "статус недействителен"
// @Security     JwtAuth
// @Router       /bindings/{uuid}/accept-reject [put]
func (h *BindingHandler) AcceptRejectBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	moderatorID := c.Get("user_id").(string)
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

// Binding godoc
// @Summary      Delete binding by UUID
// @Description  Delete binding by UUID
// @Tags         Binding
// @Produce      json
// @Param        uuid path string true "UUID"
// @Success      200  {string}  string "OK"
// @Failure      400  {string}  string "недействительный идентификатор"
// @Failure      400  {string}  string "идентификатор пользователя пустой или недействительный"
// @Failure			 401	{string}  string "отсутствует токен авторизации"
// @Failure      403  {string}  string "токен авторизации в блеклисте"
// @Failure      404  {string}  string "заявка не найдена"
// @Failure      400  {string}  string "статус не является entered"
// @Security     JwtAuth
// @Router       /bindings/{uuid} [delete]
func (h *BindingHandler) DeleteBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	userID := c.Get("user_id").(string)
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
