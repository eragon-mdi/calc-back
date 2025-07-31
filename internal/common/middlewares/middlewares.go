package middlewares

import (
	"github.com/eragon-mdi/calc-back/internal/common/configs"
	"github.com/labstack/echo/v4"
)

type Middleware interface {
	AuthToken(echo.HandlerFunc) echo.HandlerFunc
}

type customMiddleware struct {
	cfg *configs.Middlerware
}

func New(cfg *configs.Middlerware) Middleware {
	return &customMiddleware{
		cfg: cfg,
	}
}

func (m customMiddleware) AuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_ = c.Request().Header.Get("Authorization")

		// simple white list
		//if token != m.cfg.AuthToken {
		//	return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		//}

		return next(c)
	}
}
