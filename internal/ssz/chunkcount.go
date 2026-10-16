package ssz

func chunkCountContainer(v any) uint64 {
	return uint64(structValue(v).NumField())
}

func elementSize[T any]() uint64 {
	var zero T
	switch any(zero).(type) {
	case bool, uint8:
		return 1
	case uint16:
		return 2
	case uint32:
		return 4
	case uint64:
		return 8
	default:
		panic("ssz: elementSize called on a non-basic type")
	}
}
