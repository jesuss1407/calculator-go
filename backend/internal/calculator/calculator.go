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
	Add      Operation = "add"
	Subtract Operation = "subtract"
	Multiply Operation = "multiply"
	Divide   Operation = "divide"
)

// Errors returned by Calculate. Compare them with errors.Is.
var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrDivisionByZero   = errors.New("cannot divide by zero")
	ErrResultOutOfRange = errors.New("result is out of range")
)

// Calculate applies op to a and b.
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
	default:
		return 0, ErrUnknownOperation
	}

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, ErrResultOutOfRange
	}
	return result, nil
}
