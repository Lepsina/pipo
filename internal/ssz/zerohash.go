package ssz

import "math/bits"

var zeroHashes = computeZeroHashes(64)

func computeZeroHashes(n int) [][32]byte {
	hashes := make([][32]byte, n)
	for i := 1; i < n; i++ {
		hashes[i] = hash(hashes[i-1], hashes[i-1])
	}
	return hashes
}

func zeroHash(depth uint8) [32]byte {
	return zeroHashes[depth]
}

func nextPowerOfTwo(v uint64) uint64 {
	if v <= 1 {
		return 1
	}
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

func getDepth(limit uint64) uint8 {
	if limit <= 1 {
		return 0
	}
	return uint8(bits.Len64(nextPowerOfTwo(limit)) - 1)
}
