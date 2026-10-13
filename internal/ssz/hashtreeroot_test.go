package ssz

import (
	"encoding/binary"
	"testing"
)

func TestUint64HashTreeRootIsPackedValueUnchanged(t *testing.T) {
	v := Uint64(123456789)

	var want [32]byte
	binary.LittleEndian.PutUint64(want[:8], uint64(v))

	got := v.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBytes32HashTreeRootIsItselfUnchanged(t *testing.T) {
	var r Bytes32
	r[0] = 0xAB
	r[31] = 0xCD

	got := r.HashTreeRoot()
	want := [32]byte(r)

	if got != want {
		t.Fatalf("got %x, want %x (a single 32-byte chunk is its own root)", got, want)
	}
}

func TestBytes48HashTreeRoot(t *testing.T) {
	var k Bytes48
	for i := range k {
		k[i] = byte(i + 1)
	}

	var chunk0, chunk1 [32]byte
	copy(chunk0[:], k[:32])
	copy(chunk1[:], k[32:48])
	want := hash(chunk0[:], chunk1[:])

	got := k.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBytes96HashTreeRoot(t *testing.T) {
	var s Bytes96
	for i := range s {
		s[i] = byte(i + 1)
	}

	var c0, c1, c2 [32]byte
	copy(c0[:], s[:32])
	copy(c1[:], s[32:64])
	copy(c2[:], s[64:96])

	pair01 := hash(c0[:], c1[:])
	z0 := zeroHash(0)
	pair2 := hash(c2[:], z0[:])
	want := hash(pair01[:], pair2[:])

	got := s.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestUint256HashTreeRootIsItselfUnchanged(t *testing.T) {
	var v Uint256
	v[0] = 0x01
	v[31] = 0xFF

	got := v.HashTreeRoot()
	want := [32]byte(v)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBytes20HashTreeRoot(t *testing.T) {
	var a Bytes20
	for i := range a {
		a[i] = byte(i + 1)
	}

	var chunk [32]byte
	copy(chunk[:], a[:])
	want := chunk

	got := a.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBytes256HashTreeRoot(t *testing.T) {
	var l Bytes256
	for i := range l {
		l[i] = byte(i % 256)
	}

	chunks := chunksFromBytes(l[:])
	want := merklize(chunks, uint64(len(chunks)))

	got := l.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
	if len(chunks) != 8 {
		t.Fatalf("256 bytes should be exactly 8 chunks, got %d", len(chunks))
	}
}

func TestHashTreeRootByteVector(t *testing.T) {
	b := make([]byte, 64)
	b[0] = 0xFF
	b[63] = 0xAA

	var c0, c1 [32]byte
	copy(c0[:], b[:32])
	copy(c1[:], b[32:64])
	want := hash(c0[:], c1[:])

	got := HashTreeRootByteVector(b)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashTreeRootVectorOfBytes48(t *testing.T) {
	items := make([]Bytes48, 512)
	for i := range items {
		items[i][0] = byte(i % 256)
	}

	chunks := make([][32]byte, len(items))
	for i, it := range items {
		chunks[i] = it.HashTreeRoot()
	}
	want := merklize(chunks, uint64(len(chunks)))

	got := HashTreeRootVector(items)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashTreeRootVectorOfBytes32(t *testing.T) {
	for _, n := range []int{4, 6, 7} {
		items := make([]Bytes32, n)
		for i := range items {
			items[i][0] = byte(i + 1)
		}

		chunks := make([][32]byte, len(items))
		for i, it := range items {
			chunks[i] = it.HashTreeRoot()
		}
		want := merklize(chunks, uint64(len(chunks)))

		got := HashTreeRootVector(items)

		if got != want {
			t.Fatalf("n=%d: got %x, want %x", n, got, want)
		}
	}
}

func TestHashTreeRootByteList(t *testing.T) {
	b := []byte{1, 2, 3}

	var chunk [32]byte
	copy(chunk[:], b)
	want := mixInLength(chunk, uint64(len(b)))

	got := HashTreeRootByteList(b, 32)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashTreeRootByteListEmpty(t *testing.T) {
	want := mixInLength(zeroHash(0), 0)
	got := HashTreeRootByteList(nil, 32)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashTreeRootByteListMaxSize(t *testing.T) {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(i + 1)
	}

	var chunk [32]byte
	copy(chunk[:], b)
	want := mixInLength(chunk, 32)

	got := HashTreeRootByteList(b, 32)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

type twoFieldContainer struct {
	A Bytes32
	B Bytes32
}

func TestHashTreeRootContainerGeneric(t *testing.T) {
	c := twoFieldContainer{A: Bytes32{0xAA}, B: Bytes32{0xBB}}

	want := merklize([][32]byte{c.A.HashTreeRoot(), c.B.HashTreeRoot()}, 2)
	got := HashTreeRootContainer(c)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}
