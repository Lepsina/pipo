package ssz

type HashTreeRooter interface {
	HashTreeRoot() [32]byte
}

func HashTreeRootContainer(v any) [32]byte {
	val := structValue(v)
	chunks := make([][32]byte, val.NumField())
	for i := range chunks {
		hr, ok := val.Field(i).Interface().(HashTreeRooter)
		if !ok {
			panic("ssz: container field does not implement HashTreeRooter")
		}
		chunks[i] = hr.HashTreeRoot()
	}
	return merklize(chunks, chunkCountContainer(v))
}

func HashTreeRootByteVector(b []byte) [32]byte {
	chunks := chunksFromBytes(b)
	return merklize(chunks, uint64(len(chunks)))
}

func HashTreeRootVector[T HashTreeRooter](items []T) [32]byte {
	chunks := packElements(items)
	return merklize(chunks, uint64(len(chunks)))
}

func HashTreeRootByteList(b []byte, maxBytes uint64) [32]byte {
	maxChunks := (maxBytes + chunkSize - 1) / chunkSize
	root := merklize(chunksFromBytes(b), maxChunks)
	return mixInLength(root, uint64(len(b)))
}

func hashTreeRootUint64(v uint64) [32]byte {
	chunks := packBasicElements([]uint64{v})
	return merklize(chunks, uint64(len(chunks)))
}

func (v Uint64) HashTreeRoot() [32]byte {
	return hashTreeRootUint64(uint64(v))
}

func (v Uint256) HashTreeRoot() [32]byte {
	return [32]byte(v)
}

func (b Bytes20) HashTreeRoot() [32]byte {
	return HashTreeRootByteVector(b[:])
}

func (b Bytes32) HashTreeRoot() [32]byte {
	return HashTreeRootByteVector(b[:])
}

func (b Bytes48) HashTreeRoot() [32]byte {
	return HashTreeRootByteVector(b[:])
}

func (b Bytes96) HashTreeRoot() [32]byte {
	return HashTreeRootByteVector(b[:])
}

func (b Bytes256) HashTreeRoot() [32]byte {
	return HashTreeRootByteVector(b[:])
}
