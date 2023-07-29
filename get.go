package structquery

import (
	"reflect"
	"strings"
	"fmt"
	"strconv"
)

// Get retrieves the value from the object at the specified path.
// If a wildcard is used, a slice of matching values is returned.
func Get(obj interface{}, path string) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	if path == "" {
		return []interface{}{v.Interface()}, nil
	}

	pathParts := strings.Split(path, ".")
	var results []interface{}

	if err := getRecursive(v, pathParts, 0, "", &results); err != nil {
		return nil, err
	}

	return results, nil
}

// getRecursive is a helper function to implement nested wildcards.
func getRecursive(v reflect.Value, pathParts []string, pathIndex int, currPath string, results *[]interface{}) error {
	if pathIndex >= len(pathParts) {
		*results = append(*results, v.Interface())
		return nil
	}

	part := pathParts[pathIndex]

	switch v.Kind() {
	case reflect.Struct:
		if part == "*" {
			for i := 0; i < v.NumField(); i++ {
				// Avoid global matching: if the path starts with a specific field name, only consider values that match.
				if pathIndex == 0 && v.Type().Field(i).Name != pathParts[0] {
					continue
				}
				fieldVal := v.Field(i)
				if err := getRecursive(fieldVal, pathParts, pathIndex+1, currPath+v.Type().Field(i).Name+".", results); err != nil {
					return err
				}
			}
		} else {
			fieldVal := v.FieldByName(part)
			if fieldVal.IsValid() {
				if err := getRecursive(fieldVal, pathParts, pathIndex+1, currPath+part+".", results); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		val := v.MapIndex(reflect.ValueOf(part))
		if val.IsValid() {
			if err := getRecursive(val, pathParts, pathIndex+1, currPath+part+".", results); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 0 || idx >= v.Len() {
			return fmt.Errorf("invalid array index: %s", part)
		}
		val := v.Index(idx)
		if err := getRecursive(val, pathParts, pathIndex+1, currPath+strconv.Itoa(idx)+".", results); err != nil {
			return err
		}
	}
	return nil
}

