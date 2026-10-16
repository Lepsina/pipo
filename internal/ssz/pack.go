package ssz

import "encoding/binary"

func chunksFromBytes(buf []byte) [][32]byte {
	n := (len(buf) + chunkSize - 1) / chunkSize
	chunks := make([][32]byte, n)
	for i := range chunks {
		start := i * chunkSize
		end := min(start+chunkSize, len(buf))
		copy(chunks[i][:], buf[start:end])
	}
	return chunks
}

func appendBasicValue[T any](buf []byte, v T) []byte {
	switch x := any(v).(type) {
	case bool:
		if x {
			return append(buf, 1)
		}
		return append(buf, 0)
	case uint8:
		return append(buf, x)
	case uint16:
		return binary.LittleEndian.AppendUint16(buf, x)
	case uint32:
		return binary.LittleEndian.AppendUint32(buf, x)
	case uint64:
		return binary.LittleEndian.AppendUint64(buf, x)
	default:
		panic("ssz: appendBasicValue called on a non-basic type")
	}
}

func packBasicElements[T any](values []T) [][32]byte {
	buf := make([]byte, 0, len(values)*int(elementSize[T]()))
	for _, v := range values {
		buf = appendBasicValue(buf, v)
	}
	return chunksFromBytes(buf)
}

func packElements[T HashTreeRooter](values []T) [][32]byte {
	chunks := make([][32]byte, len(values))
	for i, v := range values {
		chunks[i] = v.HashTreeRoot()
	}
	return chunks
}
