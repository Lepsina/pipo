package ssz

import "testing"

type fourFieldContainer struct {
	A uint64
	B uint64
	C uint64
	D uint64
}

func TestChunkCountContainer(t *testing.T) {
	got := chunkCountContainer(fourFieldContainer{})
	if got != 4 {
		t.Fatalf("got %d, want 4", got)
	}
}

func TestChunkCountContainerViaPointer(t *testing.T) {
	got := chunkCountContainer(&fourFieldContainer{})
	if got != 4 {
		t.Fatalf("got %d, want 4", got)
	}
}
