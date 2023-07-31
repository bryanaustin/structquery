package structquery

import (
	"reflect"
	"strings"
	"github.com/ETLHero/cram"
	"github.com/hashicorp/go-multierror"
	"fmt"
)

// Add adds the provided value at the specified path. Overrides the value if target is
// a field. If target is a slice, append to slice. If target is a map and the key
// doesn't exist, create the key with value. If target is a map and the key does exist,
// override the value.
func Add(obj interface{}, path string, value interface{}) error {
	p := strings.Split(path, ".")
	last := p[len(p)-1]
	vobj := reflect.ValueOf(obj)
	v := reflect.ValueOf(value)
	if vobj.Kind() == reflect.Ptr && !vobj.IsNil() {
		vobj = vobj.Elem()
	}
	
	vs, err := Match(obj, p)
	if err != nil {
		return fmt.Errorf("matching: %w", err)
	}
	
	for i := range vs {
		if vs[i].Child == nil {
			// vs[i].Parent is a Map
			//TODO: Write a test for this above claim
			kt := vs[i].Parent.Type().Key()
			ka := reflect.New(kt)
			err := cram.Into(ka.Interface(), last)
			if err != nil {
				err = multierror.Append(err, fmt.Errorf("adding new key to map: %w", err))
				continue
			}
			if !vs[i].Parent.CanSet() {
				err = multierror.Append(err, fmt.Errorf("adding new key to map: %w", err))
				continue
			}
			vs[i].Parent.SetMapIndex(ka.Elem(), v)
			continue
		}
		if !vs[i].Child.CanSet() {
			err = multierror.Append(err, fmt.Errorf("adding new key to map: %w", err))
			continue
		}
		vs[i].Child.Set(v)
	}
	return err
}
