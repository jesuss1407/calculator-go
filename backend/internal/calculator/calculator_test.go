package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	negativeZero := math.Copysign(0, -1)

	tests := []struct {
		name    string
		op      Operation
		a, b    float64
		want    float64
		wantErr error
	}{
		// add
		{name: "add integers", op: Add, a: 2, b: 3, want: 5},
		{name: "add negatives", op: Add, a: -2, b: -3, want: -5},
		{name: "add decimals", op: Add, a: 0.5, b: 0.25, want: 0.75},
		{name: "add keeps float64 precision", op: Add, a: 0.1, b: 0.2, want: 0.30000000000000004},

		// subtract
		{name: "subtract", op: Subtract, a: 10, b: 4, want: 6},
		{name: "subtract to negative", op: Subtract, a: 4, b: 10, want: -6},

		// multiply
		{name: "multiply", op: Multiply, a: 3, b: 4, want: 12},
		{name: "multiply negatives", op: Multiply, a: -3, b: -4, want: 12},
		{name: "multiply by zero", op: Multiply, a: 3, b: 0, want: 0},

		// divide
		{name: "divide to non-integer", op: Divide, a: 10, b: 4, want: 2.5},
		{name: "divide negative", op: Divide, a: -9, b: 3, want: -3},
		{name: "divide zero by number", op: Divide, a: 0, b: 5, want: 0},

		// division by zero
		{name: "divide by zero", op: Divide, a: 1, b: 0, wantErr: ErrDivisionByZero},
		{name: "divide zero by zero", op: Divide, a: 0, b: 0, wantErr: ErrDivisionByZero},
		{name: "divide by negative zero", op: Divide, a: 1, b: negativeZero, wantErr: ErrDivisionByZero},

		// overflow
		{name: "add overflow", op: Add, a: math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrResultOutOfRange},
		{name: "multiply overflow", op: Multiply, a: 1e308, b: 10, wantErr: ErrResultOutOfRange},
		{name: "multiply negative overflow", op: Multiply, a: -1e308, b: 10, wantErr: ErrResultOutOfRange},
		{name: "divide overflow", op: Divide, a: 1e308, b: 1e-308, wantErr: ErrResultOutOfRange},

		// non-finite inputs never produce a non-finite result
		{name: "infinite input", op: Add, a: math.Inf(1), b: 1, wantErr: ErrResultOutOfRange},
		{name: "NaN input", op: Add, a: math.NaN(), b: 1, wantErr: ErrResultOutOfRange},

		// unknown operation
		{name: "unknown operation", op: "modulo", a: 1, b: 2, wantErr: ErrUnknownOperation},
		{name: "empty operation", op: "", a: 1, b: 2, wantErr: ErrUnknownOperation},
		{name: "operation is case-sensitive", op: "ADD", a: 1, b: 2, wantErr: ErrUnknownOperation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.op, tt.a, tt.b)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Calculate(%q, %v, %v) error = %v, want %v", tt.op, tt.a, tt.b, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Calculate(%q, %v, %v) = %v, want %v", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}
