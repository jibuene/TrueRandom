package randomnumber

import (
	"testing"
)

func TestNonRandomNumber(t *testing.T) {
	rng := New(NotRandom)
	rng.Generate()
	wantNumber := int64(NotRandomNumber)
	if rng.Number.Int64() != wantNumber {
		t.Errorf("Generate() = %d, want %d", rng.Number.Int64(), wantNumber)
	}
	wantBasis := "Not Random"
	if rng.Basis != wantBasis {
		t.Errorf("Generate() basis = %q, want %q", rng.Basis, wantBasis)
	}
}

func TestIsInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Simple case",
			input:    "12345",
			expected: true,
		},
		{
			name:     "With letters",
			input:    "12a34b5",
			expected: false,
		},
		{
			name:     "No digits",
			input:    "abc",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Negative number",
			input:    "-12345",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInteger(tt.input)
			if result != tt.expected {
				t.Errorf("isInteger(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNumberFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
		raises   bool
	}{
		{
			name:     "Simple case",
			input:    []string{"abc123", "def456"},
			expected: "30138989104782314344446965046",
			raises:   false,
		},
		{
			name:     "With special characters",
			input:    []string{"!@#1$%^2&*()3", "[]{}4;:'\",.<>?/5"},
			expected: "172646658740433481299052529543750197",
			raises:   false,
		},
		{
			name:     "No digits",
			input:    []string{"abc", "def"},
			expected: "107075202213222",
			raises:   false,
		},
		{
			name:     "Empty strings",
			input:    []string{"", ""},
			expected: "",
			raises:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.raises {
				assertPanic(t, func() {
					numberFromString(tt.input...)
				})
				return
			}
			result := numberFromString(tt.input...)
			if result.String() != tt.expected {
				t.Errorf("numberFromString(%v) = %q; want %q", tt.input, result.String(), tt.expected)
			}
		})
	}
}

// assertPanic checks if the provided function panics.
func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
