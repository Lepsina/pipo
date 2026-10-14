package ssz

import (
	"encoding/binary"
	"testing"
)

func TestMerklizeSingleChunkLimitOne(t *testing.T) {
	var c [32]byte
	c[0] = 0xAA
	got := merklize([][32]byte{c}, 1)
	if got != c {
		t.Fatalf("limit=1,count=1 should return the chunk unchanged")
	}
}

func TestMerklizeEmptyReturnsZeroHashAtDepth(t *testing.T) {
	got := merklize(nil, 4)
	want := zeroHash(getDepth(4))
	if got != want {
		t.Fatalf("got %x, want zeroHash(depth)=%x", got, want)
	}
}

func TestMerklizeEmptyLimitZero(t *testing.T) {
	got := merklize(nil, 0)
	if got != zeroHash(0) {
		t.Fatalf("limit=0,count=0 should return zeroHash(0)")
	}
}

func TestMerklizeTwoChunks(t *testing.T) {
	var c0, c1 [32]byte
	c0[0], c1[0] = 1, 2

	got := merklize([][32]byte{c0, c1}, 2)
	want := hash(c0, c1)
	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestMerklizeOneRealChunkPaddedToLimitFour(t *testing.T) {
	var l0 [32]byte
	l0[0] = 0xAA

	got := merklize([][32]byte{l0}, 4)

	z0 := zeroHash(0)
	z1 := zeroHash(1)
	top := hash(l0, z0)
	want := hash(top, z1)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestMerklizeOddLeftoverUsesZeroHashNotRealPadding(t *testing.T) {
	var c0, c1, c2 [32]byte
	c0[0], c1[0], c2[0] = 1, 2, 3

	got := merklize([][32]byte{c0, c1, c2}, 8)

	pair01 := hash(c0, c1)
	z0 := zeroHash(0)
	pair2 := hash(c2, z0)
	level1 := hash(pair01, pair2)
	z2 := zeroHash(2)
	want := hash(level1, z2)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestMerklizePanicsWhenCountExceedsLimit(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when count > limit")
		}
	}()
	var c0, c1 [32]byte
	merklize([][32]byte{c0, c1}, 1)
}

func TestMixInLengthMatchesHandComputedReference(t *testing.T) {
	var root [32]byte
	root[0] = 0xAB

	length := uint64(5)

	var lengthChunk [32]byte
	binary.LittleEndian.PutUint64(lengthChunk[:8], length)

	want := hash(root, lengthChunk)
	got := mixInLength(root, length)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestMixInLengthZero(t *testing.T) {
	var root [32]byte
	root[0] = 0x01

	var zeroLengthChunk [32]byte
	want := hash(root, zeroLengthChunk)
	got := mixInLength(root, 0)

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}
