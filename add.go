package structquery

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"reflect"
	"strings"
)

// Add the provided value at the specified path. Overrides the value if target is
// a field. If target is a slice, append to slice. If target is a map and the key
// doesn't exist, create the key with value. If target is a map and the key does exist,
// override the value.
func Add(obj any, path string, value any) error {
	p := strings.Split(path, ".")
	last := p[len(p)-1]
	v := reflect.ValueOf(value)
	vs, err := Match(obj, p)
	if err != nil {
		return fmt.Errorf("matching: %w", err)
	}

	for i := range vs {
		if vs[i].Child == nil {
			// vs[i].Parent is a Map
			//TODO: Write a test for this above claim
			if serr := setMap(*vs[i].Parent, last, v); serr != nil {
				err = multierror.Append(err, fmt.Errorf("adding new key to map: %w", serr))
			}
			continue
		}
		// Both parent and child exist
		if vs[i].Parent != nil {
			// Map index already exists, just overriding
			if serr := setMap(*vs[i].Parent, last, v); serr != nil {
				err = multierror.Append(err, fmt.Errorf("fallback add to new key to map: %w", serr))
			}
			continue
		}
		if vs[i].Child.Kind() == reflect.Slice {
			v = reflect.Append(*vs[i].Child, v)
		}
		vs[i].Child.Set(v)
	}
	return err
}
