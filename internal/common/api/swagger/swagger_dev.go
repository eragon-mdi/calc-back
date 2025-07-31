//go:build dev

package swagapi

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/eragon-mdi/calc-back/docs"
)

func Register(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
