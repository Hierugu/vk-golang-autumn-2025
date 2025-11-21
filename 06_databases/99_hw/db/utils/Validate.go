package utils

import (
	"fmt"
	"hw5/types"
	"strconv"
	"strings"
)

func Validate(value map[string]interface{}, columns []types.Column) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(columns))
	for _, col := range columns {
		v, exists := value[col.Name]
		dt := strings.ToLower(col.DataType)
		if !exists {
			if col.Nullable {
				result[col.Name] = nil
				continue
			} else {
				// assign zero value depending on column data type for NOT NULL columns
				if strings.HasPrefix(dt, "varchar") || dt == "text" {
					result[col.Name] = ""
					continue
				}
				if dt == "int" || strings.HasPrefix(dt, "int") {
					result[col.Name] = 0
					continue
				}
				if dt == "float" || strings.HasPrefix(dt, "float") || strings.HasPrefix(dt, "double") {
					result[col.Name] = 0.0
					continue
				}
				// fallback zero value: empty string
				result[col.Name] = ""
				continue
			}
		}

		// Handle varchar(n)
		if strings.HasPrefix(dt, "varchar") {
			// accept only string or []byte; numeric values are invalid for string columns
			switch t := v.(type) {
			case string:
				s := t
				// extract max length from varchar(n)
				if l := strings.Index(dt, "("); l != -1 {
					if r := strings.Index(dt, ")"); r != -1 && r > l+1 {
						if max, err := strconv.Atoi(dt[l+1 : r]); err == nil {
							if len(s) > max {
								return nil, fmt.Errorf("%s field exceeds maximum length of %d", col.Name, max)
							}
						}
					}
				}
				result[col.Name] = s
				continue
			case []byte:
				s := string(t)
				if l := strings.Index(dt, "("); l != -1 {
					if r := strings.Index(dt, ")"); r != -1 && r > l+1 {
						if max, err := strconv.Atoi(dt[l+1 : r]); err == nil {
							if len(s) > max {
								return nil, fmt.Errorf("%s field exceeds maximum length of %d", col.Name, max)
							}
						}
					}
				}
				result[col.Name] = s
				continue
			default:
				return nil, fmt.Errorf("field %s have invalid type", col.Name)
			}
		}

		// text -> string-like: accept only string or []byte
		if dt == "text" {
			switch t := v.(type) {
			case string:
				result[col.Name] = t
				continue
			case []byte:
				result[col.Name] = string(t)
				continue
			default:
				return nil, fmt.Errorf("field %s have invalid type", col.Name)
			}
		}

		// integer types
		if dt == "int" || strings.HasPrefix(dt, "int") {
			switch t := v.(type) {
			case int:
				result[col.Name] = t
			case int64:
				result[col.Name] = int(t)
			case float64:
				// JSON numbers are float64; verify integer value
				if t == float64(int64(t)) {
					result[col.Name] = int(t)
				} else {
					return nil, fmt.Errorf("field %s have invalid type", col.Name)
				}
			case string:
				iv, err := strconv.Atoi(t)
				if err != nil {
					return nil, fmt.Errorf("field %s have invalid type", col.Name)
				}
				result[col.Name] = iv
			case []byte:
				iv, err := strconv.Atoi(string(t))
				if err != nil {
					return nil, fmt.Errorf("field %s have invalid type", col.Name)
				}
				result[col.Name] = iv
			default:
				return nil, fmt.Errorf("field %s have invalid type", col.Name)
			}
			continue
		}

		// float types
		if dt == "float" || strings.HasPrefix(dt, "float") || strings.HasPrefix(dt, "double") {
			switch t := v.(type) {
			case float64:
				result[col.Name] = t
			case float32:
				result[col.Name] = float64(t)
			case int:
				result[col.Name] = float64(t)
			case int64:
				result[col.Name] = float64(t)
			case string:
				fv, err := strconv.ParseFloat(t, 64)
				if err != nil {
					return nil, fmt.Errorf("field %s have invalid type", col.Name)
				}
				result[col.Name] = fv
			case []byte:
				fv, err := strconv.ParseFloat(string(t), 64)
				if err != nil {
					return nil, fmt.Errorf("field %s have invalid type", col.Name)
				}
				result[col.Name] = fv
			default:
				return nil, fmt.Errorf("field %s have invalid type", col.Name)
			}
			continue
		}

		// fallback: accept only string-like or nil; otherwise invalid type
		switch v.(type) {
		case string, []byte, nil:
			result[col.Name] = v
		default:
			return nil, fmt.Errorf("field %s have invalid type", col.Name)
		}
	}
	return result, nil
}
