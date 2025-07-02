package request

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
	Tag     string      `json:"tag"`
}

// Error implements the error interface
func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", v.Field, v.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Validate validates a struct based on validation tags
func Validate(obj interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("validation can only be applied to structs")
	}

	rt := rv.Type()
	var errors ValidationErrors

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		if err := validateField(field, fieldType, validateTag); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors = append(errors, ve)
			} else {
				errors = append(errors, ValidationError{
					Field:   fieldType.Name,
					Value:   field.Interface(),
					Message: err.Error(),
					Tag:     validateTag,
				})
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateField validates a specific field
func validateField(field reflect.Value, fieldType reflect.StructField, validateTag string) error {
	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		if err := validateRule(field, fieldType, rule); err != nil {
			return err
		}
	}

	return nil
}

// validateRule validates a specific rule
func validateRule(field reflect.Value, fieldType reflect.StructField, rule string) error {
	fieldName := fieldType.Name
	fieldValue := field.Interface()

	switch {
	case rule == "required":
		if isEmpty(field) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: "field is required",
				Tag:     "required",
			}
		}

	case strings.HasPrefix(rule, "min="):
		minStr := strings.TrimPrefix(rule, "min=")
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", minStr)
		}

		if err := validateMin(field, fieldName, fieldValue, min); err != nil {
			return err
		}

	case strings.HasPrefix(rule, "max="):
		maxStr := strings.TrimPrefix(rule, "max=")
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", maxStr)
		}

		if err := validateMax(field, fieldName, fieldValue, max); err != nil {
			return err
		}

	case rule == "email":
		if field.Kind() == reflect.String {
			email := field.String()
			// Skip validation if field is empty and not required
			if email == "" {
				return nil
			}
			if !isValidEmail(email) {
				return ValidationError{
					Field:   fieldName,
					Value:   fieldValue,
					Message: "must be a valid email address",
					Tag:     "email",
				}
			}
		}

	case rule == "url":
		if field.Kind() == reflect.String {
			url := field.String()
			// Skip validation if field is empty and not required
			if url == "" {
				return nil
			}
			if !isValidURL(url) {
				return ValidationError{
					Field:   fieldName,
					Value:   fieldValue,
					Message: "must be a valid URL",
					Tag:     "url",
				}
			}
		}

	case strings.HasPrefix(rule, "regex="):
		pattern := strings.TrimPrefix(rule, "regex=")
		if field.Kind() == reflect.String {
			value := field.String()
			// Skip validation if field is empty and not required
			if value == "" {
				return nil
			}
			matched, err := regexp.MatchString(pattern, value)
			if err != nil {
				return fmt.Errorf("invalid regex pattern: %s", pattern)
			}
			if !matched {
				return ValidationError{
					Field:   fieldName,
					Value:   fieldValue,
					Message: fmt.Sprintf("must match pattern: %s", pattern),
					Tag:     "regex",
				}
			}
		}

	case rule == "numeric":
		// Skip validation if field is empty and not required
		if isEmpty(field) {
			return nil
		}
		if !isNumeric(field) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: "must be numeric",
				Tag:     "numeric",
			}
		}

	case rule == "alpha":
		if field.Kind() == reflect.String {
			value := field.String()
			// Skip validation if field is empty and not required
			if value == "" {
				return nil
			}
			if !isAlpha(value) {
				return ValidationError{
					Field:   fieldName,
					Value:   fieldValue,
					Message: "must contain only alphabetic characters",
					Tag:     "alpha",
				}
			}
		}

	case rule == "alphanumeric":
		if field.Kind() == reflect.String {
			value := field.String()
			// Skip validation if field is empty and not required
			if value == "" {
				return nil
			}
			if !isAlphanumeric(value) {
				return ValidationError{
					Field:   fieldName,
					Value:   fieldValue,
					Message: "must contain only alphanumeric characters",
					Tag:     "alphanumeric",
				}
			}
		}
	}

	return nil
}

// validateMin validates minimum constraints
func validateMin(field reflect.Value, fieldName string, fieldValue interface{}, min int) error {
	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < min {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at least %d characters long", min),
				Tag:     "min",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(min) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at least %d", min),
				Tag:     "min",
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < uint64(min) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at least %d", min),
				Tag:     "min",
			}
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < float64(min) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at least %d", min),
				Tag:     "min",
			}
		}
	case reflect.Slice, reflect.Array:
		if field.Len() < min {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must contain at least %d items", min),
				Tag:     "min",
			}
		}
	}
	return nil
}

// validateMax validates maximum constraints
func validateMax(field reflect.Value, fieldName string, fieldValue interface{}, max int) error {
	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > max {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at most %d characters long", max),
				Tag:     "max",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(max) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at most %d", max),
				Tag:     "max",
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > uint64(max) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at most %d", max),
				Tag:     "max",
			}
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > float64(max) {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must be at most %d", max),
				Tag:     "max",
			}
		}
	case reflect.Slice, reflect.Array:
		if field.Len() > max {
			return ValidationError{
				Field:   fieldName,
				Value:   fieldValue,
				Message: fmt.Sprintf("must contain at most %d items", max),
				Tag:     "max",
			}
		}
	}
	return nil
}

// isEmpty checks if a field is empty
func isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Bool:
		return !field.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	default:
		return false
	}
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// Basic email regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidURL validates URL format
func isValidURL(url string) bool {
	if url == "" {
		return false
	}

	// Basic URL regex pattern
	pattern := `^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})+(?:/.*)?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// isNumeric checks if field contains numeric value
func isNumeric(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		_, err := strconv.ParseFloat(field.String(), 64)
		return err == nil
	default:
		return false
	}
}

// isAlpha checks if string contains only alphabetic characters
func isAlpha(s string) bool {
	if s == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z]+$`, s)
	return matched
}

// isAlphanumeric checks if string contains only alphanumeric characters
func isAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, s)
	return matched
}
