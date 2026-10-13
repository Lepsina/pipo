package ssz

import (
	"encoding/hex"
	"testing"
)

func TestZeroHashBaseCaseIsLiteralZero(t *testing.T) {
	z0 := zeroHash(0)
	for _, b := range z0 {
		if b != 0 {
			t.Fatalf("zeroHash(0) should be all zero bytes, got %x", z0)
		}
	}
}

func TestZeroHashLevelsAreDistinct(t *testing.T) {
	z0 := zeroHash(0)
	z1 := zeroHash(1)
	z2 := zeroHash(2)

	if z0 == z1 || z1 == z2 || z0 == z2 {
		t.Fatalf("zeroHash levels must differ: z0=%x z1=%x z2=%x", z0, z1, z2)
	}
}

func TestZeroHashMatchesIndependentSHA256Reference(t *testing.T) {
	z1 := zeroHash(1)
	want := "f5a5fd42d16a20302798ef6ed309979b43003d2320d9f0e8ea9831a92759fb4b"
	got := hex.EncodeToString(z1[:])
	if got != want {
		t.Fatalf("zeroHash(1) = %s, want %s", got, want)
	}
}

func TestZeroHashIsRecursive(t *testing.T) {
	for depth := uint8(1); depth < 5; depth++ {
		prev := zeroHash(depth - 1)
		want := hash(prev[:], prev[:])
		got := zeroHash(depth)
		if got != want {
			t.Fatalf("zeroHash(%d) should be hash(zeroHash(%d), zeroHash(%d))", depth, depth-1, depth-1)
		}
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	cases := map[uint64]uint64{
		0: 1, 1: 1, 2: 2, 3: 4, 4: 4, 5: 8, 6: 8, 8: 8, 9: 16,
	}
	for in, want := range cases {
		if got := nextPowerOfTwo(in); got != want {
			t.Errorf("nextPowerOfTwo(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestGetDepth(t *testing.T) {
	cases := map[uint64]uint8{
		0: 0, 1: 0, 2: 1, 3: 2, 4: 2, 5: 3, 8: 3, 9: 4,
	}
	for in, want := range cases {
		if got := getDepth(in); got != want {
			t.Errorf("getDepth(%d) = %d, want %d", in, got, want)
		}
	}
}
