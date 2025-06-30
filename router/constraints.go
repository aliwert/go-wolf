package router

import (
	"regexp"
	"strconv"
	"strings"
)

// Constraint represents a parameter constraint
type Constraint func(value string) bool

// Built-in constraints
var (
	// IsNumeric checks if the value contains only digits
	IsNumeric = func(value string) bool {
		for _, r := range value {
			if r < '0' || r > '9' {
				return false
			}
		}
		return len(value) > 0
	}

	// IsAlpha checks if the value contains only letters
	IsAlpha = func(value string) bool {
		for _, r := range value {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				return false
			}
		}
		return len(value) > 0
	}

	// IsAlphaNumeric checks if the value contains only letters and digits
	IsAlphaNumeric = func(value string) bool {
		for _, r := range value {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
				return false
			}
		}
		return len(value) > 0
	}

	// IsEmail validates email format
	IsEmail = func(value string) bool {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		return emailRegex.MatchString(value)
	}

	// IsUUID validates UUID format
	IsUUID = func(value string) bool {
		uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		return uuidRegex.MatchString(value)
	}

	// IsURL validates URL format
	IsURL = func(value string) bool {
		return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
	}

	// IsSlug validates slug format (URL-friendly string)
	IsSlug = func(value string) bool {
		slugPattern := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		return slugPattern.MatchString(value)
	}

	// IsDate validates date format (YYYY-MM-DD)
	IsDate = func(value string) bool {
		datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		return datePattern.MatchString(value)
	}
)

// MinLength creates a constraint that checks minimum length
func MinLength(min int) Constraint {
	return func(value string) bool {
		return len(value) >= min
	}
}

// MaxLength creates a constraint that checks maximum length
func MaxLength(max int) Constraint {
	return func(value string) bool {
		return len(value) <= max
	}
}

// LengthRange creates a constraint that checks length range
func LengthRange(min, max int) Constraint {
	return func(value string) bool {
		l := len(value)
		return l >= min && l <= max
	}
}

// MinValue creates a constraint that checks minimum numeric value
func MinValue(min int) Constraint {
	return func(value string) bool {
		if val, err := strconv.Atoi(value); err == nil {
			return val >= min
		}
		return false
	}
}

// MaxValue creates a constraint that checks maximum numeric value
func MaxValue(max int) Constraint {
	return func(value string) bool {
		if val, err := strconv.Atoi(value); err == nil {
			return val <= max
		}
		return false
	}
}

// ValueRange creates a constraint that checks numeric value range
func ValueRange(min, max int) Constraint {
	return func(value string) bool {
		if val, err := strconv.Atoi(value); err == nil {
			return val >= min && val <= max
		}
		return false
	}
}

// OneOf creates a constraint that checks if value is one of the allowed values
func OneOf(allowed ...string) Constraint {
	allowedMap := make(map[string]bool)
	for _, v := range allowed {
		allowedMap[v] = true
	}

	return func(value string) bool {
		return allowedMap[value]
	}
}

// Regex creates a constraint that validates against a regular expression
func Regex(pattern string) Constraint {
	regex := regexp.MustCompile(pattern)
	return func(value string) bool {
		return regex.MatchString(value)
	}
}

// Custom creates a custom constraint from a function
func Custom(fn func(string) bool) Constraint {
	return fn
}

// And combines multiple constraints with AND logic
func And(constraints ...Constraint) Constraint {
	return func(value string) bool {
		for _, constraint := range constraints {
			if !constraint(value) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple constraints with OR logic
func Or(constraints ...Constraint) Constraint {
	return func(value string) bool {
		for _, constraint := range constraints {
			if constraint(value) {
				return true
			}
		}
		return false
	}
}

// Not negates a constraint
func Not(constraint Constraint) Constraint {
	return func(value string) bool {
		return !constraint(value)
	}
}
