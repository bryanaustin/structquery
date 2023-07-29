package structquery

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Remove removes a value at the specified path.
func Remove(obj interface{}, path string) error {
	pathParts := strings.Split(path, ".")

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return removeRecursive(v, pathParts, 0)
}

// removeRecursive is a helper function to implement removing values.
func removeRecursive(v reflect.Value, pathParts []string, pathIndex int) error {
	if pathIndex >= len(pathParts) {
		return fmt.Errorf("invalid path: index out of range")
	}

	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	part := pathParts[pathIndex]

	switch v.Kind() {
	case reflect.Struct:
		if pathIndex == len(pathParts)-1 {
			if part == "*" {
				errs := make([]error, 0, v.NumField())
				for i := 0; i < v.NumField(); i++ {
					field := v.Field(i)
					errs = append(errs, tryZero(field, part))
				}
				if len(errs) > 0 {
					//TODO: Return all
					return errs[0]
				}
				return nil
			} else {
				field := v.FieldByName(part)
				return tryZero(field, part)
			}
		} else {
			field := v.FieldByName(part)
			if !field.IsValid() {
				return fmt.Errorf("invalid field name: %s", part)
			}
			return removeRecursive(field, pathParts, pathIndex+1)
		}
	case reflect.Map:
		if pathIndex == len(pathParts)-1 {
			if part == "*" {
				return tryZero(v, part)
			} else {
				v.SetMapIndex(reflect.ValueOf(part), reflect.Value{})
			}
			return nil
		} else {
			val := v.MapIndex(reflect.ValueOf(part))
			if val.IsValid() {
				return removeRecursive(val, pathParts, pathIndex+1)
			} else {
				return fmt.Errorf("path not found: %s", strings.Join(pathParts[:pathIndex+1], "."))
			}
		}
	case reflect.Slice:
		if pathIndex == len(pathParts)-1 {
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= v.Len() {
				return fmt.Errorf("invalid array index: %s", part)
			}
			if part == "*" {
				return tryZero(v, part)
			}
			v.Set(appendSlice(v, index))
		} else {
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= v.Len() {
				return fmt.Errorf("invalid array index: %s", part)
			}
			val := v.Index(index)
			return removeRecursive(val, pathParts, pathIndex+1)
		}
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}

	return nil
}

// appendSlice is a helper function that appends elements of the slice excluding the given index.
func appendSlice(v reflect.Value, index int) reflect.Value {
	return reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, v.Len()))
}

// tryZero will try to zero out a value
func tryZero(v reflect.Value, name string) error {
	if !v.IsValid() {
		return fmt.Errorf("invalid value: %s", name)
	}
	if !v.CanSet() {
		return fmt.Errorf("cannot set: %s", name)
	}
	v.Set(reflect.Zero(v.Type()))
	return nil
}
