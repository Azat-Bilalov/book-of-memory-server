package app

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (a *Application) WithOptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			c.Set("user_id", "")
			return next(c)
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		err := a.redis.CheckJWTInBlacklist(c.Request().Context(), tokenString)
		if err == nil {
			return c.JSON(http.StatusForbidden, "Токен авторизации в блеклисте")
		}
		if !errors.Is(err, redis.Nil) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Неверный токен авторизации")
		}
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Неверный токен авторизации")
		}

		claims := token.Claims.(*ds.JWTClaims)

		c.Set("user_id", claims.User_id)
		c.Set("role", claims.Role)

		return next(c)
	}
}

func (a *Application) WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, "Отсутствует токен авторизации")
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		err := a.redis.CheckJWTInBlacklist(c.Request().Context(), tokenString)
		if err == nil {
			return c.JSON(http.StatusForbidden, "Токен авторизации в блеклисте")
		}
		if !errors.Is(err, redis.Nil) {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		token, err := jwt.ParseWithClaims(tokenString, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Неверный токен авторизации")
		}
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Неверный токен авторизации")
		}

		claims := token.Claims.(*ds.JWTClaims)

		c.Set("user_id", claims.User_id)
		c.Set("role", claims.Role)

		return next(c)
	}
}

func (a *Application) WithRole(requiredRole string, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role")
		log.Println("WithRole", requiredRole, role)
		if role == nil {
			return c.JSON(http.StatusUnauthorized, "Отсутствует токен авторизации")
		}
		userRole := role.(string)

		if userRole != requiredRole {
			return c.JSON(http.StatusForbidden, "Недостаточно прав для выполнения операции")
		}

		return next(c)
	}
}
