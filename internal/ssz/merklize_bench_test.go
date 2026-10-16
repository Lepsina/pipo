package ssz

import (
	"strconv"
	"testing"
)

func makeChunks(n int) [][32]byte {
	chunks := make([][32]byte, n)
	for i := range chunks {
		chunks[i][0] = byte(i)
		chunks[i][1] = byte(i >> 8)
	}
	return chunks
}

func BenchmarkHash(b *testing.B) {
	var x, y [32]byte
	x[0], y[0] = 1, 2

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x = hash(x, y)
	}
}

func BenchmarkMerklize(b *testing.B) {
	sizes := []int{2, 5, 7, 16, 512}

	for _, n := range sizes {
		chunks := makeChunks(n)
		limit := uint64(nextPowerOfTwo(uint64(n)))

		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				merklize(chunks, limit)
			}
		})
	}
}
