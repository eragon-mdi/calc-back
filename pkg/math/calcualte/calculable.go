package calculable

import (
	"errors"
	"strconv"
)

var (
	ErrUnknownOperator      = errors.New("unknown operator")
	ErrStartsWithOperator   = errors.New("expression must start with an operand")
	ErrEndsWithOperator     = errors.New("expression must end with an operand")
	ErrDoubleDot            = errors.New("invalid format: consecutive dots")
	ErrConsecutiveOperators = errors.New("invalid format: consecutive operators")
	ErrInvalidCharacter     = errors.New("unknown character in expression")
)

var operations = map[rune]func(a, b float64) float64{
	'+': func(a, b float64) float64 { return a + b },
	'-': func(a, b float64) float64 { return a - b },
	'*': func(a, b float64) float64 { return a * b },
	'/': func(a, b float64) float64 { return a / b },
}

type expression struct {
	operands  []float64
	operators []rune
}

type Calculable interface {
	GetExpression() string
	SetResult(float64)
}

func CalculateExpression(c Calculable) error {
	exp, err := parsing(c.GetExpression())
	if err != nil {
		return err
	}

	result := exp.operands[0]
	for i, operator := range exp.operators {
		operationFunc, exists := operations[operator]
		if !exists {
			return ErrUnknownOperator
		}
		result = operationFunc(result, exp.operands[i+1])
	}

	c.SetResult(result)
	return nil
}

// .
func parsing(input string) (res expression, err error) {

	// сначала операнд
	if !('0' <= input[0] && input[0] <= '9') {
		return res, ErrStartsWithOperator
	}
	// в конце также должен быть операнд
	if !('0' <= input[len(input)-1] && input[len(input)-1] <= '9') {
		return res, ErrEndsWithOperator
	}

	var buf = "" // буферная для res.operands[i]
	// var buf float64

	// парсим
	for _, char := range input {
		switch char {
		case ' ':
			continue
		// операнд
		case '.':
			if len(buf) > 0 && buf[len(buf)-1] == '.' {
				return res, ErrDoubleDot
			}
			fallthrough
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buf += string(char)
		// операции
		case '+', '-', '*', '/':
			if buf == "" || buf == "." {
				return res, ErrConsecutiveOperators
			}

			// преобразуем строку буфер в float64
			newOperand, err := strconv.ParseFloat(buf, 64)
			if err != nil {
				return res, err
			}
			buf = ""

			res.operands = append(res.operands, newOperand)
			res.operators = append(res.operators, char)
		default:
			return res, ErrInvalidCharacter
		}
	}

	// последний операнд
	newOperand, err := strconv.ParseFloat(buf, 64)
	if err != nil {
		return res, err
	}
	res.operands = append(res.operands, newOperand)

	return
}
