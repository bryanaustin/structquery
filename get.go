package structquery

import (
	"strings"
)

// Get retrieves the value from the object at the specified path.
// If a wildcard is used, a slice of matching values is returned.
func Get(obj any, path string) ([]any, error) {
	p := strings.Split(path, ".")
	vs, err := Match(obj, p)
	an := make([]any, 0, len(vs))
	for i := range vs {
		if vs[i].Child != nil {
			an = append(an, vs[i].Child.Interface())
		}
	}
	return an, err
}

