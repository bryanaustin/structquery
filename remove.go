package structquery

// Remove the provided value at the specified path. If the value of the target is the
// field of a struct, it will be zerored. If it is a pointer, it will not be nulled, but 
// indirect will be zeroed. If value of the target is in a map, it will be deleted from the map.
// If the target is an index of a slice it will be deleted from the slice while presering slice order.
// If the target is an index to an array, that value will be zeroed.
func Remove(obj interface{}, path string) error {
	//TODO: this
	return nil
}
