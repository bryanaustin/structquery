package structquery

import (
	"errors"
	"fmt"
	"github.com/ETLHero/cram"
	"reflect"
)

var (
	ErrCantSet = errors.New("cannot set")
)

func setMap(m reflect.Value, key string, value reflect.Value) error {
	kv := reflect.New(m.Type().Key())
	err := cram.Into(kv.Interface(), key)
	if err != nil {
		return fmt.Errorf("adding new key to map: %w", err)
	}
	// This gives a false negative, need a CanSetMapIndex receiver?
	// if m.CanSet() {
	// 	return ErrCantSet
	// }
	m.SetMapIndex(kv.Elem(), value)
	return nil
}
