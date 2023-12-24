package handler

import (
	"log"
	"strings"

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

// Document godoc
// @Summary Create document
// @Description Create document
// @Tags Document
// @Accept  json
// @Produce  json
// @Param title formData string true "title"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Success 200 {object} ds.Document
// @Failure 400 {string} string "неверный запрос"
// @Router /documents [post]
func (h *DocumentHandler) CreateDocument(c echo.Context) error {
	req := ds.DocumentRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(400, err.Error())
	}
	req.Image = image
	document, err := h.documentUsecase.CreateDocument(req)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, document)
}

// Document godoc
// @Summary Find active documents
// @Description Find active documents
// @Tags Document
// @Accept  json
// @Produce  json
// @Param title query string false "title"
// @Success 200 {array} ds.DocumentListResponse
// @Failure 403 {string} string "токен авторизации в блеклисте"
// @Failure 404 {string} string "документы не найдены"
// @Router /documents [get]
func (h *DocumentHandler) FindActiveDocuments(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if userID != "" && !IsValidUUID(userID) {
		c.JSON(400, "идентификатор пользователя не валидный")
	}
	title := strings.ToLower(c.QueryParam("title"))
	documents, err := h.documentUsecase.FindActiveDocuments(title, userID)
	if err != nil {
		return c.JSON(404, "документы не найдены")
	}
	return c.JSON(200, documents)
}

// Document godoc
// @Summary Find active document by uuid
// @Description Find active document by uuid
// @Tags Document
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Success 200 {object} ds.Document
// @Failure 404 {string} string "документ не найден"
// @Router /documents/{uuid} [get]
func (h *DocumentHandler) FindActiveDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	document, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, "документ не найден")
	}
	return c.JSON(200, document)
}

// Document godoc
// @Summary Update document by uuid
// @Description Update document by uuid
// @Tags Document
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Param body formData ds.DocumentRequest true "body"
// @Param image formData file true "image"
// @Success 200 {object} ds.Document
// @Failure 400 {string} string "неверный запрос"
// @Failure 404 {string} string "документ не найден"
// @Router /documents/{uuid} [put]
func (h *DocumentHandler) UpdateDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	req := ds.DocumentRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(400, err.Error())
	}
	req.Image = image
	document, err := h.documentUsecase.UpdateDocumentByUUID(uuid, req)
	if err != nil {
		return err
	}
	return c.JSON(200, document)
}

// Document godoc
// @Summary Delete document by uuid
// @Description Delete document by uuid
// @Tags Document
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Success 204
// @Failure 404 {string} string "документ не найден"
// @Router /documents/{uuid} [delete]
func (h *DocumentHandler) DeleteDocumentByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	err := h.documentUsecase.DeleteDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, err.Error())
	}
	return c.NoContent(204)
}

// Document godoc
// @Summary Add document to binding by uuid
// @Description Add document to binding by uuid
// @Tags Document
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Param body body ds.DocBinding true "body"
// @Success 204
// @Failure 400 {string} string "неверный запрос"
// @Failure 401 {string} string "отсутствует токен авторизации"
// @Failure 403 {string} string "токен авторизации в блеклисте"
// @Failure 404 {string} string "документ не найден"
// @Router /documents/{uuid}/binding [post]
func (h *DocumentHandler) AddDocumentToBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.JSON(400, "идентификатор документа пустой или неверный")
	}
	userID := c.Get("user_id").(string)
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя не валидный")
	}
	log.Println(userID)
	req := ds.DocBinding{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	_, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, "документ не найден")
	}
	err = h.documentUsecase.AddDocumentToBindingByUUID(uuid, userID, req)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.NoContent(204)
}

// Document godoc
// @Summary Remove document from binding by uuid
// @Description Remove document from binding by uuid
// @Tags Document
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Success 204
// @Failure 400 {string} string "неверный запрос"
// @Failure 404 {string} string "документ не найден"
// @Router /documents/{uuid}/binding [delete]
func (h *DocumentHandler) RemoveDocumentFromBindingByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	if !IsValidUUID(uuid) {
		return c.NoContent(404)
	}
	userID := c.Get("user_id").(string)
	if userID == "" || !IsValidUUID(userID) {
		return c.JSON(400, "идентификатор пользователя пустой или неверный")
	}
	_, err := h.documentUsecase.FindActiveDocumentByUUID(uuid)
	if err != nil {
		return c.JSON(404, "документ не найден")
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
