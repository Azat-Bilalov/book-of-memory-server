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
	documents, err := h.documentUsecase.FindActiveDocuments()
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

// todo: вынести в отдельный файл
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
