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

func setMap(m reflect.Value, key string, value any) error {
	kv := reflect.New(m.Type().Key())
	err := cram.Into(kv.Interface(), key)
	if err != nil {
		return fmt.Errorf("adding new key to map: %w", err)
	}
	if m.CanSet() {
		return ErrCantSet
	}
	m.SetMapIndex(kv.Elem(), reflect.ValueOf(value))
	return nil
}
