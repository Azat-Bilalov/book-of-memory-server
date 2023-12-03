package routes

import "github.com/labstack/echo/v4"

type Middlewares struct {
	WithAuth         echo.MiddlewareFunc
	WithOptionalAuth echo.MiddlewareFunc
	WithAdmin        echo.MiddlewareFunc
	WithUser         echo.MiddlewareFunc
}
