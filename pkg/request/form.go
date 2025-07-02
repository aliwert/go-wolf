package request

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// BindQuery binds query parameters to a struct
func BindQuery(r *http.Request, obj interface{}) error {
	values := r.URL.Query()
	if err := bindValues(values, obj, "query"); err != nil {
		return err
	}
	return Validate(obj)
}

// BindForm binds form data to a struct
func BindForm(r *http.Request, obj interface{}) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	if err := bindValues(r.Form, obj, "form"); err != nil {
		return err
	}
	return Validate(obj)
}

// BindPath binds URL path parameters to a struct
func BindPath(params map[string]string, obj interface{}) error {
	values := make(map[string][]string)
	for k, v := range params {
		values[k] = []string{v}
	}
	return bindValues(values, obj, "path")
}

// BindHeader binds HTTP headers to a struct
func BindHeader(r *http.Request, obj interface{}) error {
	values := make(map[string][]string)
	for k, v := range r.Header {
		values[strings.ToLower(k)] = v
	}
	return bindValues(values, obj, "header")
}

// BindAll binds data from multiple sources to a struct
func BindAll(r *http.Request, params map[string]string, obj interface{}) error {
	// Parse form first if applicable
	if IsForm(r) {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %w", err)
		}
	}

	// Bind path parameters
	if len(params) > 0 {
		if err := BindPath(params, obj); err != nil {
			return fmt.Errorf("failed to bind path params: %w", err)
		}
	}

	// Bind query parameters
	if err := BindQuery(r, obj); err != nil {
		return fmt.Errorf("failed to bind query params: %w", err)
	}

	// Bind form data if present
	if IsForm(r) && len(r.Form) > 0 {
		if err := BindForm(r, obj); err != nil {
			return fmt.Errorf("failed to bind form data: %w", err)
		}
	}

	// Bind headers if needed
	if err := BindHeader(r, obj); err != nil {
		return fmt.Errorf("failed to bind headers: %w", err)
	}

	return Validate(obj)
}

// bindValues binds url.Values to a struct using reflection
func bindValues(values map[string][]string, obj interface{}, tag string) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("obj must be a pointer to struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		if !field.CanSet() {
			continue
		}

		// Get tag name
		tagName := fieldType.Tag.Get(tag)
		if tagName == "" {
			tagName = strings.ToLower(fieldType.Name)
		}

		// Skip if tag is "-"
		if tagName == "-" {
			continue
		}

		// Get value from form/query
		value := values[tagName]
		if len(value) == 0 {
			continue
		}

		// Set field value based on type
		if err := setFieldValue(field, value[0]); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue sets a field value based on its type
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value: %s", value)
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
