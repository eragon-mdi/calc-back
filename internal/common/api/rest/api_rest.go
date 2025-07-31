package apirest

import (
	"github.com/eragon-mdi/calc-back/internal/common/middlewares"
	"github.com/labstack/echo/v4"
)

type Transport interface {
	GetLastCalculations(c echo.Context) error
	GetCalculationById(c echo.Context) error
	PostCalculation(c echo.Context) error
	DeleteCalcById(c echo.Context) error
	PatchCalculationById(c echo.Context) error
}

func RegisterCalculation(e *echo.Echo, t Transport, m middlewares.Middleware) {
	group := e.Group("/calculations")

	group.GET("", t.GetLastCalculations)
	group.GET("/:id", t.GetCalculationById)
	group.POST("", t.PostCalculation)
	group.DELETE("/:id", t.DeleteCalcById)
	group.PATCH("/:id", t.PatchCalculationById)
}
