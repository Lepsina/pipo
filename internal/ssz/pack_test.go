package ssz

import "testing"

func TestChunksFromBytesExactMultiple(t *testing.T) {
	buf := make([]byte, 64)
	chunks := chunksFromBytes(buf)
	if len(chunks) != 2 {
		t.Fatalf("got %d chunks, want 2", len(chunks))
	}
}

func TestChunksFromBytesPadsLastChunk(t *testing.T) {
	buf := []byte{1, 2, 3}
	chunks := chunksFromBytes(buf)
	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}
	want := [32]byte{1, 2, 3}
	if chunks[0] != want {
		t.Fatalf("got %x, want %x", chunks[0], want)
	}
}

func TestChunksFromBytesEmpty(t *testing.T) {
	chunks := chunksFromBytes(nil)
	if len(chunks) != 0 {
		t.Fatalf("got %d chunks, want 0", len(chunks))
	}
}

func TestPackBasicElementsUint64(t *testing.T) {
	chunks := packBasicElements([]uint64{1, 2, 3, 4})
	if len(chunks) != 1 {
		t.Fatalf("got %d chunks, want 1", len(chunks))
	}
	want := [32]byte{
		1, 0, 0, 0, 0, 0, 0, 0,
		2, 0, 0, 0, 0, 0, 0, 0,
		3, 0, 0, 0, 0, 0, 0, 0,
		4, 0, 0, 0, 0, 0, 0, 0,
	}
	if chunks[0] != want {
		t.Fatalf("got %x, want %x", chunks[0], want)
	}
}

func TestPackBasicElementsSpillsIntoSecondChunk(t *testing.T) {
	chunks := packBasicElements([]uint64{1, 2, 3, 4, 5})
	if len(chunks) != 2 {
		t.Fatalf("got %d chunks, want 2", len(chunks))
	}
	if chunks[1][0] != 5 {
		t.Fatalf("second chunk should start with the 5th value, got %x", chunks[1])
	}
}

func TestPackBasicElementsBool(t *testing.T) {
	chunks := packBasicElements([]bool{true, false, true})
	want := [32]byte{1, 0, 1}
	if chunks[0] != want {
		t.Fatalf("got %x, want %x", chunks[0], want)
	}
}

func TestPackBasicElementsUint8FillsOneChunk(t *testing.T) {
	values := make([]uint8, 32)
	for i := range values {
		values[i] = uint8(i)
	}
	chunks := packBasicElements(values)
	if len(chunks) != 1 {
		t.Fatalf("32 uint8 values should fit in exactly 1 chunk, got %d", len(chunks))
	}
}
