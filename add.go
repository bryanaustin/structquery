package structquery

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
)

// Add adds a new value at the specified path.
func Add(obj interface{}, path string, value interface{}) error {
	pathParts := strings.Split(path, ".")

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return addRecursive(v, pathParts, 0, value)
}

// addRecursive is a helper function to implement adding values.
func addRecursive(v reflect.Value, pathParts []string, pathIndex int, value interface{}) error {
	if pathIndex >= len(pathParts) {
		return fmt.Errorf("invalid path: index out of range")
	}

	part := pathParts[pathIndex]

	switch v.Kind() {
	case reflect.Struct:
		if pathIndex == len(pathParts)-1 {
			if part == "*" {
			} else {
				field := v.FieldByName(part)
				if !field.IsValid() {
					return fmt.Errorf("invalid field name: %s", part)
				}
				if !field.CanSet() {
					return fmt.Errorf("cannot set field: %s", part)
				}
				field.Set(reflect.ValueOf(value))
			}
		} else {
			field := v.FieldByName(part)
			if !field.IsValid() {
				return fmt.Errorf("invalid field name: %s", part)
			}
			return addRecursive(field, pathParts, pathIndex+1, value)
		}
	case reflect.Map:
		if pathIndex == len(pathParts)-1 {
			if !v.CanSet() {
				return fmt.Errorf("cannot set field: %s", part)
			}
			key := reflect.ValueOf(part)
			val := reflect.ValueOf(value)
			v.SetMapIndex(key, val)
		} else {
			val := v.MapIndex(reflect.ValueOf(part))
			if val.IsValid() {
				return addRecursive(val, pathParts, pathIndex+1, value)
			} else {
				return fmt.Errorf("path not found: %s", strings.Join(pathParts[:pathIndex+1], "."))
			}
		}
	case reflect.Slice:
		if pathIndex == len(pathParts)-1 {
			if v.Type().Elem().Kind() != reflect.ValueOf(value).Kind() {
				return fmt.Errorf("value type %s does not match slice element type %s (parent type %s)", reflect.ValueOf(value).Type(), v.Type().Elem(), reflect.ValueOf(value).Type())
			}
			newVal := reflect.Append(v, reflect.ValueOf(value))
			v.Set(newVal)
		} else {
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= v.Len() {
				return fmt.Errorf("invalid array index: %s", part)
			}
			val := v.Index(index)
			return addRecursive(val, pathParts, pathIndex+1, value)
		}
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}

	return nil
}
