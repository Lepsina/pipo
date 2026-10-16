package ssz

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Interfaces.

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal(data []byte) error
}

// fixedSizer reports whether a type is fixed-size in its SSZ encoding
// and, if so, how many bytes it always occupies. ok=false means
// variable-size (size is meaningless in that case).
type fixedSizer interface {
	FixedSize() (size uint64, ok bool)
}

// ---- Marshal ----

func copyBytes(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

func MarshalByteVector(b []byte) ([]byte, error) {
	return copyBytes(b), nil
}

func MarshalVector[T Marshaler](items []T) ([]byte, error) {
	var buf []byte
	for _, item := range items {
		b, err := item.Marshal()
		if err != nil {
			return nil, err
		}
		buf = append(buf, b...)
	}
	return buf, nil
}

func MarshalByteList(b []byte) ([]byte, error) {
	return copyBytes(b), nil
}

func MarshalContainer(v any) ([]byte, error) {
	val := structValue(v)

	type fieldBytes struct {
		b     []byte
		fixed bool
	}
	fields := make([]fieldBytes, val.NumField())
	var fixedPartLen uint64
	for i := range fields {
		m, ok := val.Field(i).Interface().(Marshaler)
		if !ok {
			panic("ssz: container field does not implement Marshaler")
		}
		b, err := m.Marshal()
		if err != nil {
			return nil, err
		}
		fs, ok := val.Field(i).Interface().(fixedSizer)
		if !ok {
			panic("ssz: container field does not implement fixedSizer")
		}
		_, fixed := fs.FixedSize()
		fields[i] = fieldBytes{b: b, fixed: fixed}
		if fixed {
			fixedPartLen += uint64(len(b))
		} else {
			fixedPartLen += offsetSize
		}
	}

	fixedBuf := make([]byte, 0, fixedPartLen)
	var variableBuf []byte
	variableOffset := fixedPartLen
	for _, f := range fields {
		if f.fixed {
			fixedBuf = append(fixedBuf, f.b...)
			continue
		}
		fixedBuf = binary.LittleEndian.AppendUint32(fixedBuf, uint32(variableOffset))
		variableBuf = append(variableBuf, f.b...)
		variableOffset += uint64(len(f.b))
	}
	return append(fixedBuf, variableBuf...), nil
}

func (v Uint64) Marshal() ([]byte, error) {
	return appendBasicValue(nil, uint64(v)), nil
}

func (v Uint256) Marshal() ([]byte, error) {
	return MarshalByteVector(v[:])
}

func (b Bytes20) Marshal() ([]byte, error) {
	return MarshalByteVector(b[:])
}

func (b Bytes32) Marshal() ([]byte, error) {
	return MarshalByteVector(b[:])
}

func (b Bytes48) Marshal() ([]byte, error) {
	return MarshalByteVector(b[:])
}

func (b Bytes96) Marshal() ([]byte, error) {
	return MarshalByteVector(b[:])
}

func (b Bytes256) Marshal() ([]byte, error) {
	return MarshalByteVector(b[:])
}

// ---- Unmarshal ----

// UnmarshalByteVector decodes a fixed-length byte vector into dst, which
// must already have the vector's length.
func UnmarshalByteVector(dst []byte, data []byte) error {
	if len(data) != len(dst) {
		return fmt.Errorf("%w: expected %d bytes, got %d", ErrLengthMismatch, len(dst), len(data))
	}
	copy(dst, data)
	return nil
}

// UnmarshalerPtr lets a generic helper call Unmarshal on each element of a
// slice without the caller needing to know the concrete pointer type.
type UnmarshalerPtr[T any] interface {
	*T
	Unmarshaler
}

// UnmarshalVector decodes Vector[T, N] with fixed-size composite elements
// into dst, which must already have N elements of elemSize bytes each.
func UnmarshalVector[T any, PT UnmarshalerPtr[T]](dst []T, data []byte, elemSize uint64) error {
	want := uint64(len(dst)) * elemSize
	if uint64(len(data)) != want {
		return fmt.Errorf("%w: expected %d bytes, got %d", ErrLengthMismatch, want, len(data))
	}
	for i := range dst {
		start := uint64(i) * elemSize
		if err := PT(&dst[i]).Unmarshal(data[start : start+elemSize]); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalByteList decodes ByteList[maxBytes] into a new slice, so the
// result never aliases data.
func UnmarshalByteList(data []byte, maxBytes uint64) ([]byte, error) {
	if uint64(len(data)) > maxBytes {
		return nil, fmt.Errorf("%w: max %d bytes, got %d", ErrListTooLong, maxBytes, len(data))
	}
	return copyBytes(data), nil
}

// UnmarshalContainer is MarshalContainer's inverse: reads the fixed part
// (inline values or 4-byte offsets) in one pass, then decodes each
// variable-size field from the trailing region using those offsets as
// boundaries. v must be a pointer to the container.
func UnmarshalContainer(v any, data []byte) error {
	val := reflect.ValueOf(v).Elem()

	type variableField struct {
		idx    int
		offset uint64
	}
	var fixedOffset uint64
	var variables []variableField
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fs, ok := field.Interface().(fixedSizer)
		if !ok {
			panic("ssz: container field does not implement fixedSizer")
		}
		size, fixed := fs.FixedSize()
		if fixed {
			if fixedOffset+size > uint64(len(data)) {
				return fmt.Errorf("%w: container truncated at field %d", ErrTruncated, i)
			}
			u, ok := field.Addr().Interface().(Unmarshaler)
			if !ok {
				panic("ssz: container field does not implement Unmarshaler")
			}
			if err := u.Unmarshal(data[fixedOffset : fixedOffset+size]); err != nil {
				return err
			}
			fixedOffset += size
			continue
		}
		if fixedOffset+offsetSize > uint64(len(data)) {
			return fmt.Errorf("%w: container truncated reading offset for field %d", ErrTruncated, i)
		}
		offset := uint64(binary.LittleEndian.Uint32(data[fixedOffset : fixedOffset+offsetSize]))
		variables = append(variables, variableField{idx: i, offset: offset})
		fixedOffset += offsetSize
	}

	if len(variables) == 0 && fixedOffset != uint64(len(data)) {
		return fmt.Errorf("%w: expected %d bytes, got %d", ErrLengthMismatch, fixedOffset, len(data))
	}
	if len(variables) > 0 && variables[0].offset != fixedOffset {
		return fmt.Errorf("%w: first offset %d does not match fixed part length %d", ErrInvalidOffset, variables[0].offset, fixedOffset)
	}

	for j, vf := range variables {
		end := uint64(len(data))
		if j+1 < len(variables) {
			end = variables[j+1].offset
		}
		if end < vf.offset || end > uint64(len(data)) {
			return fmt.Errorf("%w: at field %d", ErrInvalidOffset, vf.idx)
		}
		field := val.Field(vf.idx)
		u, ok := field.Addr().Interface().(Unmarshaler)
		if !ok {
			panic("ssz: container field does not implement Unmarshaler")
		}
		if err := u.Unmarshal(data[vf.offset:end]); err != nil {
			return err
		}
	}
	return nil
}

func (v *Uint64) Unmarshal(data []byte) error {
	if len(data) != uint64Size {
		return fmt.Errorf("%w: expected %d bytes, got %d", ErrLengthMismatch, uint64Size, len(data))
	}
	*v = Uint64(binary.LittleEndian.Uint64(data))
	return nil
}

func (v *Uint256) Unmarshal(data []byte) error {
	return UnmarshalByteVector(v[:], data)
}

func (b *Bytes20) Unmarshal(data []byte) error {
	return UnmarshalByteVector(b[:], data)
}

func (b *Bytes32) Unmarshal(data []byte) error {
	return UnmarshalByteVector(b[:], data)
}

func (b *Bytes48) Unmarshal(data []byte) error {
	return UnmarshalByteVector(b[:], data)
}

func (b *Bytes96) Unmarshal(data []byte) error {
	return UnmarshalByteVector(b[:], data)
}

func (b *Bytes256) Unmarshal(data []byte) error {
	return UnmarshalByteVector(b[:], data)
}

// ---- FixedSize ----

// FixedSizeContainer sums a container's fields' fixed sizes via reflect.
// ok=false if any field is variable-size (container is variable-size too).
func FixedSizeContainer(v any) (uint64, bool) {
	val := structValue(v)
	var total uint64
	for i := 0; i < val.NumField(); i++ {
		fs, ok := val.Field(i).Interface().(fixedSizer)
		if !ok {
			panic("ssz: container field does not implement fixedSizer")
		}
		size, fixed := fs.FixedSize()
		if !fixed {
			return 0, false
		}
		total += size
	}
	return total, true
}

func (Uint64) FixedSize() (uint64, bool) { return uint64Size, true }

func (Uint256) FixedSize() (uint64, bool) { return uint256Size, true }

func (Bytes20) FixedSize() (uint64, bool) { return bytes20Size, true }

func (Bytes32) FixedSize() (uint64, bool) { return bytes32Size, true }

func (Bytes48) FixedSize() (uint64, bool) { return bytes48Size, true }

func (Bytes96) FixedSize() (uint64, bool) { return bytes96Size, true }

func (Bytes256) FixedSize() (uint64, bool) { return bytes256Size, true }
