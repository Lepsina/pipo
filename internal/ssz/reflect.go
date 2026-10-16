package ssz

import "reflect"

// structValue returns the reflect.Value of the struct v describes,
// dereferencing one level if v is a pointer to a struct. Shared by every
// container dispatcher (chunkCountContainer, HashTreeRootContainer,
// FixedSizeContainer, MarshalContainer) since all of them accept either
// a struct value or a pointer to one.
func structValue(v any) reflect.Value {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
