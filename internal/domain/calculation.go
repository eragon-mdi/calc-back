package domain

type Calculation struct {
	ID         string
	Expression string
	Result     string
}

type CalcID struct {
	ID string
}

type CalcExpr struct {
	Expr string
}
