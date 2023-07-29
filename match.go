package structquery

import (
	"errors"
	"fmt"
	"github.com/ETLHero/cram"
	"reflect"
	"strings"
	"strconv"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrInvalidIndex = errors.New("invlaid index")
	ErrInvalidFieldName = errors.New("invlaid field name")
	ErrInvalidKey = errors.New("invlaid key")
)

// Match will find elements of the provided obj and return refrences to their values
func Match(obj interface{}, path string) ([]reflect.Value, []error) {
	v := reflect.ValueOf(obj)
	p := strings.Split(path, ".")
	return matchRecursive(v, p)
}

func matchRecursive(v reflect.Value, path []string) ([]reflect.Value, []error) {
	if len(path) < 1 {
		return []reflect.Value{v}, nil
	}

	sv := v
	r := path[0]

	for sv.Kind() == reflect.Ptr && !sv.IsNil() {
		sv = sv.Elem()
	}

	if r == "*" {
		return matchAllRecursive(sv, path)
	}
	return matchSingleRecursive(sv, path)
}

func matchSingleRecursive(v reflect.Value, path []string) ([]reflect.Value, []error) {
	r := path[0]
	nupath := path[1:]

	switch v.Kind() {
		case reflect.Struct:
			c := v.FieldByName(r)
			if !c.IsValid() {
				return nil, []error{fmt.Errorf("%w: %s", ErrInvalidFieldName, r)}
			}
			return matchRecursive(c, nupath)
		case reflect.Map:
			keytype := v.Type().Key()
			kh := reflect.New(keytype)
			if err := cram.Into(kh.Interface(), r); err != nil {
				return nil, []error{fmt.Errorf("unable to cram index %v into %s: %w", r, kh, err)}
			}
			c := v.MapIndex(kh.Elem())
			if !c.IsValid() {
				return nil, []error{fmt.Errorf("%w: %s", ErrInvalidKey, r)}
			}
			return matchRecursive(c, nupath)
		case reflect.Slice, reflect.Array:
			i, nerr := strconv.Atoi(r)
			if nerr != nil || i < 0 || i >= v.Len() {
				return nil, []error{fmt.Errorf("%w: %s lookup: %s", ErrInvalidIndex, v.Kind(), r)}
			}
			c := v.Index(i)
			return matchRecursive(c, nupath)
		default:
			return nil, []error{fmt.Errorf("%w: %s", ErrUnsupportedType, v.Type())}
	}
}

func matchAllRecursive(v reflect.Value, path []string) (svs []reflect.Value, serrs []error) {
	nupath := path[1:]

	switch v.Kind() {
		case reflect.Struct:
			n := v.NumField()
			svs = make([]reflect.Value, 0, n)
			serrs = make([]error, 0, n)
			for i := 0; i < n; i++ {
				c := v.Field(i)
				vs, errs := matchRecursive(c, nupath)
				svs = append(svs, vs...)
				serrs = append(serrs, errs...)
			}
			return
		case reflect.Map:
			n := v.Len()
			svs = make([]reflect.Value, 0, n)
			serrs = make([]error, 0, n)
			iter := v.MapRange()
			
			for iter.Next() {
				vs, errs := matchRecursive(iter.Value(), nupath)
				svs = append(svs, vs...)
				serrs = append(serrs, errs...)
			}
			return
		case reflect.Slice, reflect.Array:
			n := v.Len()
			svs = make([]reflect.Value, 0, n)
			serrs = make([]error, 0, n)
			for i := 0; i < n; i++ {
				vs, errs := matchRecursive(v.Index(i), nupath)
				svs = append(svs, vs...)
				serrs = append(serrs, errs...)
			}
			return
		default:
			return nil, []error{fmt.Errorf("%w: %s", ErrUnsupportedType, v.Type())}
	}
}

