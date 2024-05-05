package structquery

import (
	"errors"
	"fmt"
	"github.com/ETLHero/cram"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"strconv"
)

type MatchStack struct {
	Child  *reflect.Value
	Parent *reflect.Value
}

var (
	ErrUnsupportedType  = errors.New("unsupported type")
	ErrInvalidIndex     = errors.New("invlaid index")
	ErrInvalidFieldName = errors.New("invlaid field name")
	ErrInvalidKey       = errors.New("invlaid key")
	ErrNilPointer       = errors.New("nil pointer")
)

// Match will find elements of the provided obj and return refrences to their values
func Match(obj any, path []string) ([]MatchStack, error) {
	v := reflect.ValueOf(obj)
	return matchRecursive(v, path)
}

func matchRecursive(v reflect.Value, path []string) ([]MatchStack, error) {
	sv := v
	for sv.Kind() == reflect.Ptr && !sv.IsNil() {
		sv = sv.Elem()
	}

	if len(path) < 1 {
		return []MatchStack{{Child: &sv}}, nil
	}

	r := path[0]
	var vals []MatchStack
	var err error
	if r == "*" {
		vals, err = matchAllRecursive(sv, path)
	} else {
		vals, err = matchSingleRecursive(sv, path)
	}

	return vals, err
}

func matchSingleRecursive(v reflect.Value, path []string) ([]MatchStack, error) {
	r := path[0]
	nupath := path[1:]

	switch v.Kind() {
	case reflect.Struct:
		c := v.FieldByName(r)
		if !c.IsValid() {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFieldName, r)
		}
		return matchRecursive(c, nupath)
	case reflect.Map:
		keytype := v.Type().Key()
		kh := reflect.New(keytype)
		if err := cram.Into(kh.Interface(), r); err != nil {
			return nil, fmt.Errorf("unable to cram index %v into %s: %w", r, kh, err)
		}
		c := v.MapIndex(kh.Elem())
		if !c.IsValid() {
			if len(nupath) == 0 {
				return []MatchStack{{Parent: &v}}, nil
			}
			return nil, fmt.Errorf("%w: %s", ErrInvalidKey, r)
		}
		return matchRecursive(c, nupath)
	case reflect.Slice, reflect.Array:
		i, nerr := strconv.Atoi(r)
		if nerr != nil || i < 0 || i >= v.Len() {
			return nil, fmt.Errorf("%w: %s lookup: %s", ErrInvalidIndex, v.Kind(), r)
		}
		c := v.Index(i)
		return matchRecursive(c, nupath)
	case reflect.Pointer:
		if v.IsNil() {
			return nil, ErrNilPointer
		}
		return matchRecursive(v.Elem(), nupath)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, v.Type())
	}
}

func matchAllRecursive(v reflect.Value, path []string) (svs []MatchStack, serrs error) {
	nupath := path[1:]

	switch v.Kind() {
	case reflect.Struct:
		n := v.NumField()
		for i := 0; i < n; i++ {
			c := v.Field(i)
			vs, errs := matchRecursive(c, nupath)
			svs = append(svs, vs...)
			if errs != nil {
				serrs = multierror.Append(serrs, errs)
			}
		}
		return
	case reflect.Map:
		iter := v.MapRange()

		for iter.Next() {
			vs, errs := matchRecursive(iter.Value(), nupath)
			svs = append(svs, vs...)
			if errs != nil {
				serrs = multierror.Append(serrs, errs)
			}
		}
		return
	case reflect.Slice, reflect.Array:
		n := v.Len()
		for i := 0; i < n; i++ {
			vs, errs := matchRecursive(v.Index(i), nupath)
			svs = append(svs, vs...)
			if errs != nil {
				serrs = multierror.Append(serrs, errs)
			}
		}
		return
	case reflect.Pointer:
		if v.IsNil() {
			return nil, ErrNilPointer
		}
		return matchRecursive(v.Elem(), nupath)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, v.Type())
	}
}
