package ssz

import "encoding/binary"

func merklize(chunks [][32]byte, limit uint64) [32]byte {
	nodes := make([][32]byte, len(chunks))
	copy(nodes, chunks)
	return merklizeInPlace(nodes, limit)
}

func merklizeInPlace(chunks [][32]byte, limit uint64) [32]byte {
	count := uint64(len(chunks))
	if count > limit {
		panic("ssz: chunk count exceeds limit")
	}

	if limit <= 1 {
		if count == 1 {
			return chunks[0]
		}
		return zeroHash(0)
	}

	depth := getDepth(limit)
	if count == 0 {
		return zeroHash(depth)
	}

	levelLen := len(chunks)
	for level := uint8(0); level < depth; level++ {
		for i := 0; i < levelLen; i += 2 {
			left := chunks[i]
			var right [32]byte
			if i+1 < levelLen {
				right = chunks[i+1]
			} else {
				right = zeroHash(level)
			}
			chunks[i/2] = hash(left, right)
		}
		levelLen = (levelLen + 1) / 2
	}
	return chunks[0]
}

func mixInLength(root [32]byte, length uint64) [32]byte {
	var lengthChunk [32]byte
	binary.LittleEndian.PutUint64(lengthChunk[:8], length)
	return hash(root, lengthChunk)
}
