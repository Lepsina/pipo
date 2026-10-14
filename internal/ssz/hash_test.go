package ssz

import (
	"crypto/sha256"
	"testing"
)

func TestHashMatchesSHA256(t *testing.T) {
	var a, b [32]byte
	a[0], a[1], a[2] = 1, 2, 3
	b[0], b[1], b[2] = 4, 5, 6

	got := hash(a, b)
	want := sha256.Sum256(append(append([]byte{}, a[:]...), b[:]...))

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashDoesNotMutateInputs(t *testing.T) {
	var a, b [32]byte
	a[0], a[1], a[2] = 1, 2, 3
	aCopy := a
	b[0], b[1], b[2] = 4, 5, 6

	hash(a, b)

	if a != aCopy {
		t.Fatalf("hash mutated its first argument: got %v, want %v", a, aCopy)
	}
}
