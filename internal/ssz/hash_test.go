package ssz

import (
	"crypto/sha256"
	"testing"
)

func TestHashMatchesSHA256(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{4, 5, 6}

	got := hash(a, b)
	want := sha256.Sum256(append(append([]byte{}, a...), b...))

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestHashDoesNotMutateInputs(t *testing.T) {
	a := make([]byte, 3, 10)
	a[0], a[1], a[2] = 1, 2, 3
	aCopy := append([]byte{}, a...)
	b := []byte{4, 5, 6}

	hash(a, b)

	if string(a) != string(aCopy) {
		t.Fatalf("hash mutated its first argument: got %v, want %v", a, aCopy)
	}
}
