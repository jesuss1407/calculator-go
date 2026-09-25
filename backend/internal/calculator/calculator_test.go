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

		// power
		{name: "power", op: Power, a: 2, b: 10, want: 1024},
		{name: "power with negative exponent", op: Power, a: 2, b: -1, want: 0.5},
		{name: "power of negative base with integer exponent", op: Power, a: -2, b: 3, want: -8},
		{name: "power with fractional exponent", op: Power, a: 9, b: 0.5, want: 3},
		{name: "zero to the power of zero", op: Power, a: 0, b: 0, want: 1},
		{name: "zero to a negative power", op: Power, a: 0, b: -1, wantErr: ErrDivisionByZero},
		{name: "negative base with fractional exponent", op: Power, a: -4, b: 0.5, wantErr: ErrNotRealNumber},
		{name: "cube root of negative is not supported", op: Power, a: -8, b: 1.0 / 3, wantErr: ErrNotRealNumber},
		{name: "power overflow", op: Power, a: 10, b: 400, wantErr: ErrResultOutOfRange},

		// square root (unary: b is ignored)
		{name: "sqrt", op: Sqrt, a: 16, want: 4},
		{name: "sqrt of non-square", op: Sqrt, a: 2, want: 1.4142135623730951},
		{name: "sqrt of zero", op: Sqrt, a: 0, want: 0},
		{name: "sqrt ignores b", op: Sqrt, a: 9, b: 123, want: 3},
		{name: "sqrt of negative", op: Sqrt, a: -4, wantErr: ErrNotRealNumber},

		// percentage: a% of b
		{name: "percentage", op: Percentage, a: 10, b: 200, want: 20},
		{name: "percentage with decimal percent", op: Percentage, a: 12.5, b: 64, want: 8},
		{name: "percentage above 100", op: Percentage, a: 150, b: 20, want: 30},
		{name: "negative percentage", op: Percentage, a: -10, b: 50, want: -5},
		{name: "zero percent", op: Percentage, a: 0, b: 5, want: 0},
		{name: "percentage keeps float64 precision", op: Percentage, a: 7, b: 300, want: 21.000000000000004},
		{name: "percentage has no intermediate overflow", op: Percentage, a: 200, b: 1e307, want: 2e307},
		{name: "percentage overflow", op: Percentage, a: 1e308, b: 1e308, wantErr: ErrResultOutOfRange},

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

func TestIsValid(t *testing.T) {
	for _, op := range []Operation{Add, Subtract, Multiply, Divide, Power, Sqrt, Percentage} {
		if !op.IsValid() {
			t.Errorf("Operation(%q).IsValid() = false, want true", op)
		}
		// IsValid and Calculate each list the operations; make sure they agree.
		if _, err := Calculate(op, 4, 2); errors.Is(err, ErrUnknownOperation) {
			t.Errorf("Calculate(%q, 4, 2) = ErrUnknownOperation, but IsValid says it is supported", op)
		}
	}

	for _, op := range []Operation{"modulo", "", "ADD"} {
		if op.IsValid() {
			t.Errorf("Operation(%q).IsValid() = true, want false", op)
		}
	}
}

func TestIsUnary(t *testing.T) {
	tests := []struct {
		op   Operation
		want bool
	}{
		{op: Sqrt, want: true},
		{op: Add, want: false},
		{op: Power, want: false},
		{op: Percentage, want: false},
		{op: "modulo", want: false},
	}

	for _, tt := range tests {
		if got := tt.op.IsUnary(); got != tt.want {
			t.Errorf("Operation(%q).IsUnary() = %v, want %v", tt.op, got, tt.want)
		}
	}
}
