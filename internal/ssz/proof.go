package ssz

import "math/bits"

func ComputeMerkleRoot(leaf [32]byte, branch [][32]byte, depth uint64, index uint64) [32]byte {
	value := leaf
	for i := uint64(0); i < depth; i++ {
		sibling := branch[i]
		if (index>>i)&1 == 1 {
			value = hash(sibling[:], value[:])
		} else {
			value = hash(value[:], sibling[:])
		}
	}
	return value
}

func VerifyMerkleBranch(leaf [32]byte, branch [][32]byte, depth uint64, index uint64, root [32]byte) bool {
	return ComputeMerkleRoot(leaf, branch, depth, index) == root
}

func VerifyMerkleBranchAtGindex(leaf [32]byte, branch [][32]byte, gindex uint64, root [32]byte) bool {
	depth := floorLog2(gindex)
	index := subtreeIndex(gindex, depth)
	return VerifyMerkleBranch(leaf, branch, depth, index, root)
}

func floorLog2(x uint64) uint64 {
	return uint64(bits.Len64(x) - 1)
}

func subtreeIndex(gindex uint64, depth uint64) uint64 {
	return gindex - (1 << depth)
}
