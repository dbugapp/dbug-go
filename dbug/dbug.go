package dbug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

var endpoint = "http://127.0.0.1:53821"

// SetEndpoint allows setting a custom endpoint.
func SetEndpoint(url string) {
	endpoint = url
}

// Send serializes and sends the payload to the debug app.
func Send(payload interface{}) {
	jsonBytes, err := stringify(payload)
	if err != nil {
		jsonBytes, _ = json.MarshalIndent(map[string]string{
			"error":  "Serialization failed",
			"reason": err.Error(),
		}, "", "  ")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 500 * time.Millisecond}
	_, _ = client.Do(req)
}

func stringify(data interface{}) ([]byte, error) {
	safe, err := sanitize(data, map[uintptr]bool{})
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(safe, "", "  ")
}

func sanitize(data interface{}, seen map[uintptr]bool) (interface{}, error) {
	if marshaler, ok := data.(json.Marshaler); ok {
		bytes, err := marshaler.MarshalJSON()
		if err == nil {
			var resultInterface interface{}
			if json.Unmarshal(bytes, &resultInterface) == nil {
				return resultInterface, nil
			}
		}
	}

	if data == nil {
		return nil, nil
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil, nil
		}
		ptr := v.Pointer()
		if seen[ptr] {
			return "[circular]", nil
		}

		seen[ptr] = true
		defer delete(seen, ptr)

		return sanitize(v.Elem().Interface(), seen)

	case reflect.Map:
		result := map[string]interface{}{}
		for _, key := range v.MapKeys() {
			val, err := sanitize(v.MapIndex(key).Interface(), seen)
			if err != nil {
				return nil, err
			}
			result[fmt.Sprintf("%v", key.Interface())] = val
		}
		return result, nil

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			val, err := sanitize(v.Index(i).Interface(), seen)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil

	case reflect.Struct:
		result := map[string]interface{}{}
		result["go_type"] = v.Type().Name()

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.PkgPath != "" {
				continue
			}

			fieldValue := v.Field(i)
			if fieldValue.Type() == reflect.TypeOf(json.RawMessage{}) {
				var rawValue interface{}
				jsonBytes := fieldValue.Interface().(json.RawMessage)
				if json.Unmarshal(jsonBytes, &rawValue) == nil {
					result[field.Name] = rawValue
				} else {
					result[field.Name] = string(jsonBytes)
				}
				continue
			}

			val, err := sanitize(fieldValue.Interface(), seen)
			if err != nil {
				return nil, err
			}
			result[field.Name] = val
		}
		return result, nil

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return fmt.Sprintf("[%s]", v.Kind().String()), nil

	default:
		return data, nil
	}
}

// SendTestable exposes stringified output for testing purposes.
func SendTestable(payload interface{}) ([]byte, error) {
	return stringify(payload)
}
