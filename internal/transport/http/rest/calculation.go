package resttransport

import (
	"net/http"

	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//go:generate mockery --name=Service --with-expecter --output=./mocks --exported
type Service interface {
	GetLastCalculations() ([]domain.Calculation, error)
	GetCalculationById(domain.CalcID) (domain.Calculation, error)
	CreateCalculation(domain.CalcExpr) (domain.Calculation, error)
	DeleteCalcById(domain.CalcID) error
	UpdateCalculationById(domain.Calculation) (domain.Calculation, error)
}

const (
	paramID              = "id"
	logErrInvalidUUID    = "invalid UUID"
	logErrInvalidBodyReq = "invalid body request"
)

// GetLastCalculations godoc
// @Summary      Получить крайние 10 вычислений
// @Description  Возвращает массив с информацией о вычислениях
// @Tags         calculations
// @Accept       json
// @Produce      json
// @Success      200 {array} CalcResponse
// @Failure 	 400 {object} ErrorResponse
// @Failure 	 404 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /calculations [get]
func (t transport) GetLastCalculations(c echo.Context) error {
	calcs, err := t.s.GetLastCalculations()
	if err != nil {
		t.l.Error("transport.GetLastCalculations failed to get calculations", "cause", err)
		return httpErrHandler(err)
	}

	t.l.Info("transport.GetLastCalculations calculations getted successfully", "res", calcs)

	return c.JSON(http.StatusOK, calcsResponse(calcs))
}

// GetCalculationById godoc
// @Summary      Получить вычисление по индетификатору
// @Description  Возвращает информацию о вычислении. Индентификатор - строковой тип UUID
// @Tags         calculations
// @Accept       json
// @Produce      json
// @Param        id path string true "Индентификатор"
// @Success      200 {object} CalcResponse
// @Failure 	 400 {object} ErrorResponse
// @Failure 	 404 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /calculations/{id} [get]
func (t transport) GetCalculationById(c echo.Context) error {
	idStr := c.Param(paramID)

	if err := uuid.Validate(idStr); err != nil {
		t.l.Error("transport.GetCalculationById", logErrInvalidUUID, "cause", err)
		return echo.NewHTTPError(http.StatusBadRequest, errRespBadIdParam)
	}

	calc, err := t.s.GetCalculationById(calcId(idStr))
	if err != nil {
		t.l.Error("transport.GetCalculationById failed to get calculation", "cause", err)
		return httpErrHandler(err)
	}

	t.l.Info("transport.GetCalculationById calculation getted successfully", "res", calc)

	return c.JSON(http.StatusOK, calcResponse(calc))
}

// PostCalculation godoc
// @Summary      Добавить новое вычисление
// @Description  Создает новое вычисление на основе выражения
// @Tags         calculations
// @Accept       json
// @Produce      json
// @Param        request body CalcRequest true "Данные для вычисления"
// @Success      201 {object} CalcResponse
// @Failure 	 400 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /calculations [post]
func (t transport) PostCalculation(c echo.Context) error {
	var calcReq CalcRequest
	if err := c.Bind(&calcReq); err != nil {
		t.l.Error("transport.PostCalculation", logErrInvalidBodyReq, "cause", err)
		return echo.NewHTTPError(http.StatusBadRequest, errRespBadRequest)
	}

	calc, err := t.s.CreateCalculation(calcReq.CalcExpr())
	if err != nil {
		t.l.Error("transport.PostCalculation failed post calculation", "cause", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errRespInternal)
	}

	t.l.Info("transport.PostCalculation calculation created successfully", "res", calc)

	return c.JSON(http.StatusCreated, calcResponse(calc))
}

// DeleteCalcById godoc
// @Summary      Удаляет вычисление
// @Description  Удаляет вычисление по индентификатору. Индентификатор - строковой тип UUID
// @Tags         calculations
// @Accept       json
// @Produce      json
// @Param        id path string true "Индентификатор"
// @Success      204
// @Failure 	 400 {object} ErrorResponse
// @Failure 	 404 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /calculations/{id} [delete]
func (t transport) DeleteCalcById(c echo.Context) error {
	idStr := c.Param(paramID)

	if err := uuid.Validate(idStr); err != nil {
		t.l.Error("transport.DeleteCalcById ", logErrInvalidUUID, " cause ", err, " id ", idStr)
		return echo.NewHTTPError(http.StatusBadRequest, errRespBadIdParam)
	}

	if err := t.s.DeleteCalcById(calcId(idStr)); err != nil {
		t.l.Error("transport.DeleteCalcById failed to delete calculation", "cause", err)
		return httpErrHandler(err)
	}

	t.l.Info("transport.DeleteCalcById calculation deleted successfully")

	return c.JSON(http.StatusNoContent, nil)
}

// PatchCalculationById godoc
// @Summary      Изменить вычисление
// @Description  Изменяет выражение для расчёта и результат существующего выражения по индентификатору. Индентификатор - строковой тип UUID
// @Tags         calculations
// @Accept       json
// @Produce      json
// @Param        id path string true "Индентификатор"
// @Success      200 {object} CalcResponse
// @Failure 	 400 {object} ErrorResponse
// @Failure 	 404 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /calculations/{id} [patch]
func (t transport) PatchCalculationById(c echo.Context) error {
	idStr := c.Param(paramID)

	if err := uuid.Validate(idStr); err != nil {
		t.l.Error("transport.PatchCalculationById", logErrInvalidUUID, "cause", err)
		return echo.NewHTTPError(http.StatusBadRequest, errRespBadIdParam)
	}

	var calcReq CalcRequest
	if err := c.Bind(&calcReq); err != nil {
		t.l.Error("transport.PatchCalculationById", logErrInvalidBodyReq, "cause", err)
		return echo.NewHTTPError(http.StatusBadRequest, errRespBadRequest)
	}

	calc, err := t.s.UpdateCalculationById(calcReq.Calculation(idStr))
	if err != nil {
		t.l.Error("transport.PatchCalculationById failed to update calculation", "cause", err)
		return httpErrHandler(err)
	}

	t.l.Info("transport.PatchCalculationById calculation updated successfully", "res", calc)

	return c.JSON(http.StatusOK, calcResponse(calc))
}
