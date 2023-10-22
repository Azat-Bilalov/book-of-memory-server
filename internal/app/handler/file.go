package handler

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/labstack/echo/v4"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) FindFile(c echo.Context) error {
	bucket := c.Param("bucket")
	filename := c.Param("filename")

	contentBytes, contentType, err := config.ReadObject(bucket, filename)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	return c.Blob(200, contentType, contentBytes)
}
