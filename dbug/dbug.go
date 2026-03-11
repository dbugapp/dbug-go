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

// Go serializes and sends all payloads to the Dbug app in a single request.
// Multiple payloads are sent as a JSON array and appear stacked in one cell.
func Go(payloads ...interface{}) {
	if len(payloads) == 0 {
		return
	}

	var jsonBytes []byte
	var err error

	if len(payloads) == 1 {
		jsonBytes, err = stringify(payloads[0])
	} else {
		sanitized := make([]interface{}, len(payloads))
		for i, payload := range payloads {
			sanitized[i], err = sanitize(payload, map[uintptr]bool{})
			if err != nil {
				sanitized[i] = map[string]string{
					"_error_":  "Serialization failed during Send",
					"_reason_": err.Error(),
				}
				err = nil
			}
		}
		jsonBytes, err = json.MarshalIndent(sanitized, "", "  ")
	}

	if err != nil {
		jsonBytes, _ = json.MarshalIndent(map[string]string{
			"_error_":  "Serialization failed during Send",
			"_reason_": err.Error(),
		}, "", "  ")
	}

	client := http.Client{Timeout: 500 * time.Millisecond}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Printf("Dbug: Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("Dbug: Error sending request: %v\n", err)
	}
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
		structType := v.Type()
		structTypeName := structType.String()

		for i := 0; i < v.NumField(); i++ {
			field := structType.Field(i)
			fieldName := field.Name
			prefixedFieldName := fmt.Sprintf("%s.%s", structTypeName, fieldName)

			if field.PkgPath == "" {
				fieldValue := v.Field(i)

				if fieldValue.Type() == reflect.TypeOf(json.RawMessage{}) {
					var rawValue interface{}
					jsonBytes := fieldValue.Interface().(json.RawMessage)
					if json.Unmarshal(jsonBytes, &rawValue) == nil {
						result[prefixedFieldName] = rawValue
					} else {
						result[prefixedFieldName] = string(jsonBytes)
					}
					continue
				}

				val, err := sanitize(fieldValue.Interface(), seen)
				if err != nil {
					return nil, fmt.Errorf("failed to sanitize field %s: %w", prefixedFieldName, err)
				}
				result[prefixedFieldName] = val
			} else {
				// Field is unexported (private)
				// Consistently format private fields as "FieldName [FieldTypeString]"
				privateFieldValue := fmt.Sprintf("%s [%s]", field.Name, field.Type.String())
				result[prefixedFieldName] = privateFieldValue
			}
		}
		return result, nil

	case reflect.Chan:
		chanType := v.Type()
		typeString := chanType.String()
		chanDetails := map[string]interface{}{}

		if v.IsNil() {
			chanDetails["is_nil"] = true
			chanDetails["element_type"] = chanType.Elem().String()
			chanDetails["direction"] = chanDirToString(chanType.ChanDir())
			chanDetails["capacity"] = 0
			chanDetails["length"] = 0
		} else {
			chanDetails["is_nil"] = false
			chanDetails["element_type"] = chanType.Elem().String()
			chanDetails["direction"] = chanDirToString(chanType.ChanDir())
			chanDetails["capacity"] = v.Cap()
			chanDetails["length"] = v.Len()
		}

		result := map[string]interface{}{
			typeString: chanDetails,
		}
		return result, nil

	case reflect.UnsafePointer:
		if v.IsNil() {
			return "[UnsafePointer (nil)]", nil
		}
		return fmt.Sprintf("[UnsafePointer 0x%x]", v.Pointer()), nil

	case reflect.Func:
		fnType := v.Type()
		signature := fnType.String()

		numIn := fnType.NumIn()
		inTypes := make([]string, numIn)
		for i := 0; i < numIn; i++ {
			inTypes[i] = fnType.In(i).String()
		}

		numOut := fnType.NumOut()
		outTypes := make([]string, numOut)
		for i := 0; i < numOut; i++ {
			outTypes[i] = fnType.Out(i).String()
		}

		funcDetails := map[string]interface{}{
			"input_types":  inTypes,
			"output_types": outTypes,
			"is_variadic":  fnType.IsVariadic(),
		}

		result := map[string]interface{}{
			signature: funcDetails,
		}
		return result, nil

	default:
		return data, nil
	}
}

func chanDirToString(dir reflect.ChanDir) string {
	switch dir {
	case reflect.SendDir:
		return "send-only"
	case reflect.RecvDir:
		return "receive-only"
	case reflect.BothDir:
		return "send-receive"
	default:
		return "unknown"
	}
}

// GoTestable handles a single payload for testing.
func GoTestable(payload interface{}) ([]byte, error) {
	return stringify(payload)
}
