package structquery

// Get retrieves the value from the object at the specified path.
// If a wildcard is used, a slice of matching values is returned.
func Get(obj any, path string) ([]any, error) {
	vs, err := Match(obj, path)
	an := make([]any, len(vs))
	for i := range vs {
		an[i] = vs[i].Interface()
	}
	return an, err
}

