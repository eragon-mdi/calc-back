package service

import (
	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

//go:generate mockery --name=Repository --with-expecter --output=./mocks --exported
type Repository interface {
	GetCalculations(int) ([]domain.Calculation, error)
	GetCalculation(string) (domain.Calculation, error)
	DeleteCalculation(string) error
	SaveTask(domain.Calculation) (domain.Calculation, error)
	UpdateTaskInfo(domain.Calculation) (domain.Calculation, error)
}

const Max_Calcs = 10

func (s service) GetLastCalculations() ([]domain.Calculation, error) {
	calcs, err := s.r.GetCalculations(Max_Calcs)
	if err != nil {
		return nil, errors.Wrap(err, "service: failed to get last calcs")
	}

	if len(calcs) == 0 {
		return nil, domain.ErrNotFound
	}

	return calcs, nil
}

func (s service) GetCalculationById(id domain.CalcID) (domain.Calculation, error) {
	task, err := s.r.GetCalculation(id.ID)
	if err != nil {
		return domain.Calculation{}, errors.Wrap(err, "service: failed to get calc")
	}

	return task, nil
}

func (s service) CreateCalculation(expr domain.CalcExpr) (domain.Calculation, error) {
	calc, err := calculate(domain.Calculation{Expression: expr.Expr})
	if err != nil {
		return domain.Calculation{}, domain.ErrValidation
	}

	calc.ID = uuid.NewString()

	calc, err = s.r.SaveTask(calc)
	if err != nil {
		return domain.Calculation{}, errors.Wrap(err, "service: failed to save new calc")
	}

	return calc, nil
}

func (s service) DeleteCalcById(id domain.CalcID) error {
	if err := s.r.DeleteCalculation(id.ID); err != nil {
		return errors.Wrap(err, "service: failed to delete calc")
	}

	return nil
}

func (s service) UpdateCalculationById(calc domain.Calculation) (domain.Calculation, error) {
	calc, err := calculate(calc)
	if err != nil {
		return domain.Calculation{}, domain.ErrValidation
	}

	newCalc, err := s.r.UpdateTaskInfo(calc)
	if err != nil {
		return domain.Calculation{}, errors.Wrap(err, "service: failed to update calc")
	}

	return newCalc, nil
}
