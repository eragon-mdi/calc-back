package service

import (
	"strconv"

	"github.com/eragon-mdi/calc-back/internal/domain"
	calculable "github.com/eragon-mdi/calc-back/pkg/math/calcualte"
)

type calc struct {
	domain.Calculation
}

func (c calc) GetExpression() string {
	return c.Expression
}

func (c *calc) SetResult(res float64) {
	c.Result = strconv.FormatFloat(res, 'f', -1, 64)
}

func calculate(c domain.Calculation) (domain.Calculation, error) {
	newC := calc{Calculation: c}
	if err := calculable.CalculateExpression(&newC); err != nil {
		return domain.Calculation{}, err
	}

	return newC.Calculation, nil
}
