package resttransport

import (
	"errors"
	"net/http"

	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/labstack/echo/v4"
)

var (
	errRespNotFound   = ErrorResponse{"No content found"}
	errRespInternal   = ErrorResponse{"Service is currently unavailable. Please try again later."}
	errRespBadIdParam = ErrorResponse{"Bad id param"}
	errRespBadRequest = ErrorResponse{"Bad request"}
	errRespValidation = ErrorResponse{"Failed on validation"}
)

func httpErrHandler(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, errRespNotFound)
	case errors.Is(err, domain.ErrValidation):
		return echo.NewHTTPError(http.StatusBadRequest, errRespValidation)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, errRespInternal)
	}
}
