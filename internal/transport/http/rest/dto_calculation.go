package resttransport

import (
	"github.com/eragon-mdi/calc-back/internal/domain"
)

type CalcRequest struct {
	Expression string `json:"expression" example:"2+3/2"`
}

type CalcResponse struct {
	ID         string `json:"id" example:"a8098c1a-f86e-11da-bd1a-00112444be1e"`
	Expression string `json:"expression" example:"2+3/2"`
	Result     string `json:"result" example:"3.5"`
}

type ErrorResponse struct {
	Message string `json:"error" example:"<cause>"`
}

func calcId(id string) domain.CalcID {
	return domain.CalcID{
		ID: id,
	}
}

func (c CalcRequest) CalcExpr() domain.CalcExpr {
	return domain.CalcExpr{
		Expr: c.Expression,
	}
}

func (c CalcRequest) Calculation(id string) domain.Calculation {
	return domain.Calculation{
		ID:         id,
		Expression: c.Expression,
	}
}

func calcResponse(c domain.Calculation) CalcResponse {
	return CalcResponse{
		ID:         c.ID,
		Expression: c.Expression,
		Result:     c.Result,
	}
}

func calcsResponse(cs []domain.Calculation) []CalcResponse {
	res := make([]CalcResponse, 0, len(cs))
	for _, c := range cs {
		res = append(res, calcResponse(c))
	}
	return res
}
