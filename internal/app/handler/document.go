package handler

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	documentUsecase *usecase.DocumentUsecase
}

func NewDocumentHandler(documentUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{documentUsecase}
}

func (h *DocumentHandler) CreateDocument(c echo.Context) error {
	req := ds.DocumentCreateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	document, err := h.documentUsecase.CreateDocument(req)
	if err != nil {
		return err
	}
	return c.JSON(201, document)
}

func (h *DocumentHandler) FindActiveDocuments(c echo.Context) error {
	title := c.QueryParam("title")
	documents, err := h.documentUsecase.FindActiveDocuments(title)
	if err != nil {
		return err
	}
	return c.JSON(200, documents)
}

func (h *DocumentHandler) FindActiveDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	document, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return err
	}
	if document == nil {
		return c.NoContent(404)
	}
	return c.JSON(200, document)
}

func (h *DocumentHandler) UpdateDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	req := ds.DocumentUpdateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	document, err := h.documentUsecase.UpdateDocumentByUUID(uuid, req)
	if err != nil {
		return err
	}
	return c.JSON(200, document)
}

func (h *DocumentHandler) DeleteDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	err := h.documentUsecase.DeleteDocumentByUUID(uuid)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (h *DocumentHandler) FindDocumentInBinding(c echo.Context) error {
	bindingUUID := c.Param("uuid-binding")
	if !IsValidUUID(bindingUUID) {
		return c.NoContent(404)
	}
	documentUUID := c.Param("uuid-document")
	if !IsValidUUID(documentUUID) {
		return c.NoContent(404)
	}
	document, err := h.documentUsecase.FindDocumentInBinding(bindingUUID, documentUUID)
	if err != nil {
		return c.JSON(404, err.Error())
	}
	return c.JSON(200, document)
}

func (h *DocumentHandler) AddDocumentToBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "user id is empty or invalid")
	}
	req := ds.DocBindingRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	_, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, err.Error())
	}
	err = h.documentUsecase.AddDocumentToBindingByUUID(uuid, userID, req)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.NoContent(204)
}

func (h *DocumentHandler) RemoveDocumentFromBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	userID := c.Request().Header.Get("x-user-id")
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "user id is empty or invalid")
	}
	_, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, err.Error())
	}
	err = h.documentUsecase.RemoveDocumentFromBindingByUUID(uuid, userID)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.NoContent(204)
}

// todo: вынести в отдельный файл
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
