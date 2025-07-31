package api

import (
	apirest "github.com/eragon-mdi/calc-back/internal/common/api/rest"
	swagapi "github.com/eragon-mdi/calc-back/internal/common/api/swagger"
	"github.com/eragon-mdi/calc-back/internal/common/middlewares"
	"github.com/labstack/echo/v4"
)

type Transport interface {
	apirest.Transport
}

func RegisterRoutes(e *echo.Echo, t Transport, m middlewares.Middleware) {
	e.Use(middlewares.CORS(e))

	swagapi.Register(e)

	// group := e.Group("/v1", m.AuthToken)
	apirest.RegisterCalculation(e, t, m)
}
