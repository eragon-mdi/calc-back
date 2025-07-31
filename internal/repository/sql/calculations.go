package sqlrepo

import (
	"database/sql"

	"github.com/eragon-mdi/calc-back/internal/domain"
	"github.com/go-faster/errors"
)

const (
	ErrFailedQuery        = "repo: failed query"
	ErrFailedExec         = "repo: failed exec"
	ErrFailedScan         = "repo: failed to scan row"
	ErrFailedAffectedRows = "repo: failed to get number of affected rows"
	ErrFailedStartTX      = "repo: failed to start tx"
	ErrFailedCommitTX     = "repo: failed to commit tx"
	ErrFailedRollbackTX   = "repo: failed rollback tx"
)

func (r sqlRepo) GetCalculations(maxCount int) (calcs []domain.Calculation, err error) {
	rows, err := r.s.Query(getCalcsWithMax, maxCount)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedQuery)
	}
	defer rows.Close()

	for rows.Next() {
		calc := domain.Calculation{}
		if err := rows.Scan(&calc.ID, &calc.Expression, &calc.Result); err != nil {
			return nil, errors.Wrap(err, ErrFailedScan)
		}

		calcs = append(calcs, calc)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ErrFailedQuery)
	}

	return calcs, nil
}

func (r sqlRepo) GetCalculation(id string) (domain.Calculation, error) {
	var calc = domain.Calculation{}

	row := r.s.QueryRow(getCalcById, id)

	if err := row.Scan(&calc.ID, &calc.Expression, &calc.Result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Calculation{}, domain.ErrNotFound
		}
		return domain.Calculation{}, errors.Wrap(err, ErrFailedScan)
	}

	return calc, nil
}

func (r sqlRepo) DeleteCalculation(id string) error {
	res, err := r.s.Exec(deleteCalcById, id)
	if err != nil {
		return errors.Wrap(err, ErrFailedExec)
	}

	c, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, ErrFailedAffectedRows)
	}

	if c == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r sqlRepo) SaveTask(calc domain.Calculation) (domain.Calculation, error) {
	var newCalc = domain.Calculation{}

	row := r.s.QueryRow(insertCalc, calc.ID, calc.Expression, calc.Result)

	if err := row.Scan(&newCalc.ID, &newCalc.Expression, &newCalc.Result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Calculation{}, domain.ErrNotFound
		}
		return domain.Calculation{}, errors.Wrap(err, ErrFailedScan)
	}

	return newCalc, nil
}

func (r sqlRepo) UpdateTaskInfo(calc domain.Calculation) (domain.Calculation, error) {
	var updatedCalc = domain.Calculation{}

	row := r.s.QueryRow(updateCalc, calc.ID, calc.Expression, calc.Result)

	if err := row.Scan(&updatedCalc.ID, &updatedCalc.Expression, &updatedCalc.Result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Calculation{}, domain.ErrNotFound
		}
		return domain.Calculation{}, errors.Wrap(err, ErrFailedScan)
	}

	return updatedCalc, nil
}
