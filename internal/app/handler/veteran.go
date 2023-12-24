package handler

import (
	"strings"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
)

type VeteranHandler struct {
	veteranUsecase *usecase.VeteranUsecase
}

func NewVeteranHandler(veteranUsecase *usecase.VeteranUsecase) *VeteranHandler {
	return &VeteranHandler{veteranUsecase}
}

// Veteran godoc
// @Summary Create veteran
// @Description Create veteran
// @Tags Veteran
// @Accept  json
// @Produce  json
// @Param name formData string true "name"
// @Param description formData string true "description"
// @Param image formData file true "image"
// @Success 200 {object} ds.
// @Failure 400 {string} string "неверный запрос"
// @Router /veterans [post]
func (h *VeteranHandler) CreateVeteran(c echo.Context) error {
	req := ds.VeteranRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(400, err.Error())
	}
	req.Image = image
	veteran, err := h.veteranUsecase.CreateVeteran(req)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, veteran)
}

// Veteran godoc
// @Summary Find veterans
// @Description Find veterans
// @Tags Veteran
// @Accept  json
// @Produce  json
// @Param name query string false "name"
// @Success 200 {array} ds.
// @Failure 403 {string} string "токен авторизации в блеклисте"
// @Failure 404 {string} string "документы не найдены"
// @Router /veterans [get]
func (h *VeteranHandler) FindVeterans(c echo.Context) error {
	name := strings.ToLower(c.QueryParam("name"))
	veterans, err := h.veteranUsecase.FindVeterans(name)
	if err != nil {
		return c.JSON(404, "ветераны не найдены")
	}
	return c.JSON(200, veterans)
}
