package ssz

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestUint64Marshal(t *testing.T) {
	v := Uint64(0x0102030405060708)
	got, err := v.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	if string(got) != string(want) {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestUint64Unmarshal(t *testing.T) {
	data := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	var v Uint64
	if err := v.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != Uint64(0x0102030405060708) {
		t.Fatalf("got %x, want %x", v, Uint64(0x0102030405060708))
	}
}

func TestUint64UnmarshalWrongLengthErrors(t *testing.T) {
	var v Uint64
	if err := v.Unmarshal([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for wrong-length input")
	}
}

func TestUint64FixedSize(t *testing.T) {
	size, ok := Uint64(0).FixedSize()
	if !ok || size != 8 {
		t.Fatalf("got (%d,%v), want (8,true)", size, ok)
	}
}

func TestBytes32MarshalUnmarshalRoundTrip(t *testing.T) {
	var r Bytes32
	for i := range r {
		r[i] = byte(i)
	}

	data, err := r.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 32 {
		t.Fatalf("got len %d, want 32", len(data))
	}

	var got Bytes32
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != r {
		t.Fatalf("got %x, want %x", got, r)
	}
}

func TestBytes32MarshalDoesNotAliasCaller(t *testing.T) {
	var r Bytes32
	r[0] = 0xAA

	data, _ := r.Marshal()
	data[0] = 0xFF

	if r[0] != 0xAA {
		t.Fatal("mutating the returned slice must not mutate the source Bytes32")
	}
}

func TestBytes32UnmarshalWrongLengthErrors(t *testing.T) {
	var r Bytes32
	if err := r.Unmarshal([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for wrong-length input")
	}
}

func TestBytes32FixedSize(t *testing.T) {
	size, ok := Bytes32{}.FixedSize()
	if !ok || size != 32 {
		t.Fatalf("got (%d,%v), want (32,true)", size, ok)
	}
}

func TestBytes48MarshalUnmarshalRoundTrip(t *testing.T) {
	var k Bytes48
	for i := range k {
		k[i] = byte(i)
	}
	data, err := k.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got Bytes48
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != k {
		t.Fatalf("got %x, want %x", got, k)
	}
}

func TestBytes48FixedSize(t *testing.T) {
	size, ok := Bytes48{}.FixedSize()
	if !ok || size != 48 {
		t.Fatalf("got (%d,%v), want (48,true)", size, ok)
	}
}

func TestBytes96MarshalUnmarshalRoundTrip(t *testing.T) {
	var s Bytes96
	for i := range s {
		s[i] = byte(i)
	}
	data, err := s.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got Bytes96
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != s {
		t.Fatalf("got %x, want %x", got, s)
	}
}

func TestBytes96FixedSize(t *testing.T) {
	size, ok := Bytes96{}.FixedSize()
	if !ok || size != 96 {
		t.Fatalf("got (%d,%v), want (96,true)", size, ok)
	}
}

func TestUint256MarshalUnmarshalRoundTrip(t *testing.T) {
	var u Uint256
	u[0] = 0xAA
	data, err := u.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got Uint256
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != u {
		t.Fatalf("got %x, want %x", got, u)
	}
}

func TestUint256FixedSize(t *testing.T) {
	size, ok := Uint256{}.FixedSize()
	if !ok || size != 32 {
		t.Fatalf("got (%d,%v), want (32,true)", size, ok)
	}
}

func TestBytes20MarshalUnmarshalRoundTrip(t *testing.T) {
	var a Bytes20
	for i := range a {
		a[i] = byte(i)
	}
	data, err := a.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got Bytes20
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != a {
		t.Fatalf("got %x, want %x", got, a)
	}
}

func TestBytes20FixedSize(t *testing.T) {
	size, ok := Bytes20{}.FixedSize()
	if !ok || size != 20 {
		t.Fatalf("got (%d,%v), want (20,true)", size, ok)
	}
}

func TestBytes256MarshalUnmarshalRoundTrip(t *testing.T) {
	var l Bytes256
	l[0], l[255] = 0x01, 0x02
	data, err := l.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got Bytes256
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != l {
		t.Fatalf("got %x, want %x", got, l)
	}
}

func TestBytes256FixedSize(t *testing.T) {
	size, ok := Bytes256{}.FixedSize()
	if !ok || size != 256 {
		t.Fatalf("got (%d,%v), want (256,true)", size, ok)
	}
}

func TestMarshalByteVectorDoesNotAliasCaller(t *testing.T) {
	b := []byte{0xAA, 0xBB}
	data, err := MarshalByteVector(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data[0] = 0xFF
	if b[0] != 0xAA {
		t.Fatal("mutating the returned slice must not mutate the source bytes")
	}
}

func TestUnmarshalByteVectorWrongLengthErrors(t *testing.T) {
	dst := make([]byte, 64)
	if err := UnmarshalByteVector(dst, make([]byte, 63)); err == nil {
		t.Fatal("expected error for wrong-length input")
	}
}

func TestMarshalUnmarshalVectorRoundTrip(t *testing.T) {
	items := []Bytes32{{0x01}, {0x02}, {0x03}}

	data, err := MarshalVector(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 3*32 {
		t.Fatalf("got len %d, want %d", len(data), 3*32)
	}

	got := make([]Bytes32, 3)
	if err := UnmarshalVector(got, data, 32); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, items) {
		t.Fatalf("got %x, want %x", got, items)
	}
}

func TestUnmarshalVectorWrongLengthErrors(t *testing.T) {
	got := make([]Bytes32, 2)
	if err := UnmarshalVector(got, []byte{1, 2, 3}, 32); err == nil {
		t.Fatal("expected error for wrong-length input")
	}
}

func TestMarshalUnmarshalByteListRoundTrip(t *testing.T) {
	b := []byte{1, 2, 3}

	data, err := MarshalByteList(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := UnmarshalByteList(data, 32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(got, b) {
		t.Fatalf("got %x, want %x", got, b)
	}
}

func TestMarshalByteListDoesNotAliasCaller(t *testing.T) {
	b := []byte{0xAA}
	data, _ := MarshalByteList(b)
	data[0] = 0xFF
	if b[0] != 0xAA {
		t.Fatal("mutating the returned slice must not mutate the source bytes")
	}
}

func TestUnmarshalByteListDoesNotAliasInput(t *testing.T) {
	data := []byte{0xAA}
	got, err := UnmarshalByteList(data, 32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data[0] = 0xFF
	if got[0] != 0xAA {
		t.Fatal("mutating the input must not mutate the decoded list")
	}
}

func TestUnmarshalByteListAtMaxLengthSucceeds(t *testing.T) {
	if _, err := UnmarshalByteList(make([]byte, 32), 32); err != nil {
		t.Fatalf("unexpected error at max length: %v", err)
	}
}

func TestUnmarshalByteListTooLongErrors(t *testing.T) {
	if _, err := UnmarshalByteList(make([]byte, 33), 32); err == nil {
		t.Fatal("expected error for a list longer than its max")
	}
}

// testByteList is ByteList[8], a variable-size field for container tests.
type testByteList []byte

func (l testByteList) HashTreeRoot() [32]byte { return HashTreeRootByteList(l, 8) }

func (l testByteList) Marshal() ([]byte, error) { return MarshalByteList(l) }

func (l *testByteList) Unmarshal(data []byte) error {
	b, err := UnmarshalByteList(data, 8)
	if err != nil {
		return err
	}
	*l = b
	return nil
}

func (testByteList) FixedSize() (uint64, bool) { return 0, false }

type mixedContainer struct {
	A Bytes32
	B testByteList
	C Uint64
}

// mixedContainerFixedPartLen = A (32) + offset for B (4) + C (8).
const mixedContainerFixedPartLen = 32 + 4 + 8

func TestFixedSizeContainer(t *testing.T) {
	size, ok := FixedSizeContainer(twoFieldContainer{})
	if !ok || size != 64 {
		t.Fatalf("got (%d,%v), want (64,true)", size, ok)
	}
	if _, ok := FixedSizeContainer(mixedContainer{}); ok {
		t.Fatal("container with a variable-size field should be variable-size")
	}
}

func TestMarshalContainerOffsetTable(t *testing.T) {
	c := mixedContainer{A: Bytes32{0x01}, B: testByteList{0xAA, 0xBB, 0xCC}, C: 7}

	data, err := MarshalContainer(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantLen := mixedContainerFixedPartLen + len(c.B)
	if len(data) != wantLen {
		t.Fatalf("got len %d, want %d", len(data), wantLen)
	}

	gotOffset := binary.LittleEndian.Uint32(data[32:36])
	if gotOffset != mixedContainerFixedPartLen {
		t.Fatalf("got offset %d, want %d", gotOffset, mixedContainerFixedPartLen)
	}

	if !bytes.Equal(data[mixedContainerFixedPartLen:], c.B) {
		t.Fatalf("got variable part %x, want %x", data[mixedContainerFixedPartLen:], c.B)
	}
}

func TestMarshalUnmarshalContainerRoundTrip(t *testing.T) {
	c := mixedContainer{A: Bytes32{0x01}, B: testByteList{0xAA, 0xBB}, C: 7}

	data, err := MarshalContainer(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got mixedContainer
	if err := UnmarshalContainer(&got, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, c) {
		t.Fatalf("got %+v, want %+v", got, c)
	}
}

func TestUnmarshalContainerBadOffsetErrors(t *testing.T) {
	c := mixedContainer{A: Bytes32{0x01}, B: testByteList{0xAA}, C: 7}
	data, _ := MarshalContainer(c)

	// Corrupt the offset to point past the end of the buffer.
	binary.LittleEndian.PutUint32(data[32:36], uint32(len(data)+100))

	var got mixedContainer
	if err := UnmarshalContainer(&got, data); err == nil {
		t.Fatal("expected error for out-of-range offset")
	}
}
