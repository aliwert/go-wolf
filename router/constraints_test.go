package router

import (
	"testing"
)

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"999", true},
		{"", false},
		{"abc", false},
		{"12a", false},
		{"a12", false},
		{"12.5", false},
		{"-123", false},
		{"+123", false},
	}

	for _, test := range tests {
		result := IsNumeric(test.input)
		if result != test.expected {
			t.Errorf("IsNumeric(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"ABC", true},
		{"AbC", true},
		{"", false},
		{"123", false},
		{"abc123", false},
		{"ab-c", false},
		{"ab_c", false},
		{"ab c", false},
	}

	for _, test := range tests {
		result := IsAlpha(test.input)
		if result != test.expected {
			t.Errorf("IsAlpha(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC123", true},
		{"AbC123", true},
		{"abc", true},
		{"123", true},
		{"", false},
		{"abc-123", false},
		{"abc_123", false},
		{"abc 123", false},
		{"abc.123", false},
	}

	for _, test := range tests {
		result := IsAlphaNumeric(test.input)
		if result != test.expected {
			t.Errorf("IsAlphaNumeric(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.org", true},
		{"123@example.com", true},
		{"", false},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"test@domain", false},
		{"test.example.com", false},
		{"test@@example.com", false},
	}

	for _, test := range tests {
		result := IsEmail(test.input)
		if result != test.expected {
			t.Errorf("IsEmail(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"6ba7b811-9dad-11d1-80b4-00c04fd430c8", true},
		{"", false},
		{"invalid", false},
		{"550e8400-e29b-41d4-a716", false},
		{"550e8400-e29b-41d4-a716-446655440000-extra", false},
		{"550e8400e29b41d4a716446655440000", false},
		{"550e8400-e29b-41d4-a716-44665544000g", false},
	}

	for _, test := range tests {
		result := IsUUID(test.input)
		if result != test.expected {
			t.Errorf("IsUUID(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"https://www.example.com/path", true},
		{"http://localhost:8080", true},
		{"", false},
		{"example.com", false},
		{"ftp://example.com", false},
		{"//example.com", false},
		{"http//example.com", false},
	}

	for _, test := range tests {
		result := IsURL(test.input)
		if result != test.expected {
			t.Errorf("IsURL(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello-world", true},
		{"my-blog-post", true},
		{"test123", true},
		{"a", true},
		{"", false},
		{"Hello-World", false},
		{"hello_world", false},
		{"hello world", false},
		{"hello--world", false},
		{"-hello", false},
		{"hello-", false},
	}

	for _, test := range tests {
		result := IsSlug(test.input)
		if result != test.expected {
			t.Errorf("IsSlug(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2023-12-31", true},
		{"2000-01-01", true},
		{"1999-12-25", true},
		{"", false},
		{"2023-13-01", true}, // Note: This only validates format, not actual date validity
		{"2023-12-32", true}, // Note: This only validates format, not actual date validity
		{"23-12-31", false},
		{"2023/12/31", false},
		{"2023-12-1", false},
		{"2023-1-31", false},
	}

	for _, test := range tests {
		result := IsDate(test.input)
		if result != test.expected {
			t.Errorf("IsDate(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestMinLength(t *testing.T) {
	constraint := MinLength(5)

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"hello world", true},
		{"test", false},
		{"", false},
		{"exactly5", true},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("MinLength(5)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestMaxLength(t *testing.T) {
	constraint := MaxLength(5)

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"test", true},
		{"", true},
		{"12345", true},
		{"toolong", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("MaxLength(5)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestLengthRange(t *testing.T) {
	constraint := LengthRange(3, 8)

	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"test", true},
		{"12345678", true},
		{"", false},
		{"ab", false},
		{"too long string", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("LengthRange(3, 8)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestMinValue(t *testing.T) {
	constraint := MinValue(10)

	tests := []struct {
		input    string
		expected bool
	}{
		{"10", true},
		{"15", true},
		{"100", true},
		{"9", false},
		{"0", false},
		{"abc", false},
		{"", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("MinValue(10)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestMaxValue(t *testing.T) {
	constraint := MaxValue(100)

	tests := []struct {
		input    string
		expected bool
	}{
		{"100", true},
		{"50", true},
		{"0", true},
		{"101", false},
		{"200", false},
		{"abc", false},
		{"", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("MaxValue(100)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestValueRange(t *testing.T) {
	constraint := ValueRange(10, 100)

	tests := []struct {
		input    string
		expected bool
	}{
		{"10", true},
		{"50", true},
		{"100", true},
		{"9", false},
		{"101", false},
		{"abc", false},
		{"", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("ValueRange(10, 100)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestOneOf(t *testing.T) {
	constraint := OneOf("red", "green", "blue")

	tests := []struct {
		input    string
		expected bool
	}{
		{"red", true},
		{"green", true},
		{"blue", true},
		{"yellow", false},
		{"", false},
		{"Red", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("OneOf(red, green, blue)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestRegex(t *testing.T) {
	constraint := Regex(`^[a-z]+$`)

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"world", true},
		{"", false}, // Empty string should not match ^[a-z]+$ (+ means one or more)
		{"Hello", false},
		{"hello123", false},
		{"hello-world", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Regex(^[a-z]+$)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestCustom(t *testing.T) {
	customFn := func(value string) bool {
		return len(value) > 0 && value[0] == 'A'
	}
	constraint := Custom(customFn)

	tests := []struct {
		input    string
		expected bool
	}{
		{"Apple", true},
		{"Awesome", true},
		{"", false},
		{"banana", false},
		{"123", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Custom(starts with A)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestAnd(t *testing.T) {
	constraint := And(MinLength(3), IsAlpha)

	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", true},
		{"hello", true},
		{"ab", false},     // Too short
		{"abc123", false}, // Not alpha
		{"", false},       // Both conditions fail
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("And(MinLength(3), IsAlpha)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestOr(t *testing.T) {
	constraint := Or(IsNumeric, IsAlpha)

	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"abc", true},
		{"ABC", true},
		{"", false},
		{"abc123", false},
		{"ab-c", false},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Or(IsNumeric, IsAlpha)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestNot(t *testing.T) {
	constraint := Not(IsNumeric)

	tests := []struct {
		input    string
		expected bool
	}{
		{"123", false},
		{"abc", true},
		{"", true},
		{"abc123", true},
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Not(IsNumeric)(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestComplexConstraintCombinations(t *testing.T) {
	// Test complex combination: (IsAlpha OR IsNumeric) AND MinLength(2) AND MaxLength(10)
	constraint := And(
		Or(IsAlpha, IsNumeric),
		MinLength(2),
		MaxLength(10),
	)

	tests := []struct {
		input    string
		expected bool
	}{
		{"ab", true},
		{"12", true},
		{"abcdefghij", true},   // exactly 10 chars
		{"a", false},           // too short
		{"abcdefghijk", false}, // too long
		{"ab1", false},         // mixed alpha-numeric
		{"", false},            // empty
	}

	for _, test := range tests {
		result := constraint(test.input)
		if result != test.expected {
			t.Errorf("Complex constraint(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkIsNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsNumeric("123456789")
	}
}

func BenchmarkIsAlpha(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsAlpha("abcdefghijklmnop")
	}
}

func BenchmarkIsEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsEmail("test@example.com")
	}
}

func BenchmarkRegexConstraint(b *testing.B) {
	constraint := Regex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	for i := 0; i < b.N; i++ {
		constraint("test@example.com")
	}
}

func BenchmarkComplexConstraint(b *testing.B) {
	constraint := And(
		Or(IsAlpha, IsNumeric),
		MinLength(2),
		MaxLength(50),
	)
	for i := 0; i < b.N; i++ {
		constraint("validstring")
	}
}
