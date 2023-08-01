package structquery

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Replace replaces a value at the specified path.
func Replace(obj interface{}, path string, value interface{}) (interface{}, error) {
	pathParts := strings.Split(path, ".")

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return replaceRecursive(v, pathParts, 0, value)
}

// replaceRecursive is a helper function to implement replacing values.
func replaceRecursive(v reflect.Value, pathParts []string, pathIndex int, value interface{}) (interface{}, error) {
	if pathIndex >= len(pathParts) {
		return nil, fmt.Errorf("invalid path: index out of range")
	}

	part := pathParts[pathIndex]

	switch v.Kind() {
	case reflect.Struct:
		if pathIndex == len(pathParts)-1 {
			field := v.FieldByName(part)
			if !field.IsValid() {
				return nil, fmt.Errorf("invalid field name: %s", part)
			}
			if !field.CanSet() {
				return nil, fmt.Errorf("cannot set field: %s", part)
			}
			previous := field.Interface()
			field.Set(reflect.ValueOf(value))
			return previous, nil
		} else {
			field := v.FieldByName(part)
			if !field.IsValid() {
				return nil, fmt.Errorf("invalid field name: %s", part)
			}
			return replaceRecursive(field, pathParts, pathIndex+1, value)
		}
	case reflect.Map:
		if pathIndex == len(pathParts)-1 {
			key := reflect.ValueOf(part)
			val := v.MapIndex(key)
			if !val.IsValid() {
				return nil, fmt.Errorf("key not found: %s", part)
			}
			previous := val.Interface()
			v.SetMapIndex(key, reflect.ValueOf(value))
			return previous, nil
		} else {
			val := v.MapIndex(reflect.ValueOf(part))
			if !val.IsValid() {
				return nil, fmt.Errorf("key not found: %s", part)
			}
			return replaceRecursive(val, pathParts, pathIndex+1, value)
		}
	case reflect.Slice:
		index, err := strconv.Atoi(part)
		if err != nil || index < 0 || index >= v.Len() {
			return nil, fmt.Errorf("invalid array index: %s", part)
		}
		val := v.Index(index)
		if pathIndex == len(pathParts)-1 {
			previous := val.Interface()
			val.Set(reflect.ValueOf(value))
			return previous, nil
		} else {
			return replaceRecursive(val, pathParts, pathIndex+1, value)
		}
	default:
		return nil, fmt.Errorf("unsupported type: %s", v.Type())
	}
}
