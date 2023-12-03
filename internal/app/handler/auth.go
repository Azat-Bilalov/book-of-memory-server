package handler

import (
	"fmt"
	"strings"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase,
	}
}

// Auth godoc
// @Summary      Login
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 			 request body ds.LoginRequest true "request"
// @Success      200  {object}  ds.LoginResponse
// @Failure      400  {string}  string "некорректный запрос"
// @Failure      400  {string}  string "неверный логин или пароль"
// @Router       /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	req := ds.LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "некорректный запрос")
	}
	fmt.Println(req)
	res, err := h.authUsecase.Login(req.Email, req.Passwd)
	if err != nil {
		return c.JSON(400, "неверный логин или пароль")
	}
	return c.JSON(200, res)
}

// Auth godoc
// @Summary      Register
// @Description  Register
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body ds.RegisterRequest true "request"
// @Success      200  {object}  ds.RegisterResponse
// @Failure      400  {string}  string "некорректный запрос"
// @Failure      400  {string}  string "пользователь с таким email уже существует"
// @Router       /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req := ds.RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	res, err := h.authUsecase.Register(req.FirstName, req.LastName, req.Email, req.Passwd)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, res)
}

// Auth godoc
// @Summary      Logout
// @Description  Logout
// @Tags         Auth
// @Produce      json
// @Success      200  {string}  string "выход выполнен успешно"
// @Failure      400  {string}  string "неверный токен авторизации"
// @Failure      401 {string}  string "неверный токен авторизации"
// @Failure      403 {string}  string "токен авторизации в блеклисте"
// @Security     JwtAuth
// @Router       /logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(400, "неверный токен авторизации")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	err := h.authUsecase.Logout(c.Request().Context(), tokenString)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, "выход выполнен успешно")
}
