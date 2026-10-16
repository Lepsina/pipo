package ssz

import "errors"

var (
	// ErrLengthMismatch means the input is not exactly the length this
	// type's SSZ encoding requires.
	ErrLengthMismatch = errors.New("ssz: length mismatch")

	// ErrTruncated means the input ran out of bytes before a container
	// could read a field or a field's offset.
	ErrTruncated = errors.New("ssz: truncated data")

	// ErrListTooLong means a byte list's length exceeds its declared
	// maximum.
	ErrListTooLong = errors.New("ssz: list exceeds max length")

	// ErrInvalidOffset means a container's offset table is internally
	// inconsistent: the first offset doesn't match the fixed-part
	// length, or a field's region isn't within bounds.
	ErrInvalidOffset = errors.New("ssz: invalid offset table")
)
