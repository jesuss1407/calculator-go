// Package calculator implements the arithmetic behind the calculator API.
// It has no knowledge of HTTP or JSON.
package calculator

import (
	"errors"
	"math"
)

// Operation identifies an arithmetic operation.
type Operation string

// Supported operations.
const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"
	Sqrt       Operation = "sqrt"
	Percentage Operation = "percentage"
)

// Errors returned by Calculate. Compare them with errors.Is.
var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrDivisionByZero   = errors.New("cannot divide by zero")
	ErrNotRealNumber    = errors.New("result is not a real number")
	ErrResultOutOfRange = errors.New("result is out of range")
)

// IsUnary reports whether op uses only the first operand.
func (op Operation) IsUnary() bool {
	return op == Sqrt
}

// Calculate applies op to a and b. Unary operations ignore b.
// On success the result is always a finite number; on error it is 0.
func Calculate(op Operation, a, b float64) (float64, error) {
	var result float64

	switch op {
	case Add:
		result = a + b
	case Subtract:
		result = a - b
	case Multiply:
		result = a * b
	case Divide:
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result = a / b
	case Power:
		switch {
		case a == 0 && b < 0:
			return 0, ErrDivisionByZero // 0 to a negative power is 1 / 0
		case a < 0 && b != math.Trunc(b):
			return 0, ErrNotRealNumber // e.g. (-4)^0.5 has no real value
		}
		result = math.Pow(a, b)
	case Sqrt:
		if a < 0 {
			return 0, ErrNotRealNumber
		}
		result = math.Sqrt(a)
	case Percentage:
		// a% of b. Dividing first also keeps a*b from overflowing when the result itself fits.
		result = (a / 100) * b
	default:
		return 0, ErrUnknownOperation
	}

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, ErrResultOutOfRange
	}
	return result, nil
}
