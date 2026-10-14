package ssz

import "testing"

func TestVerifyMerkleBranchDepthTwoAllIndexOrders(t *testing.T) {
	var leaf, s0, s1 [32]byte
	leaf[0], s0[0], s1[0] = 0xAA, 0x01, 0x02

	cases := []struct {
		index uint64
		want  [32]byte
	}{
		{0, hashLevel(hash(leaf, s0), s1)},
		{1, hashLevel(hash(s0, leaf), s1)},
		{2, hashLevelSwap(hash(leaf, s0), s1)},
		{3, hashLevelSwap(hash(s0, leaf), s1)},
	}

	for _, c := range cases {
		got := ComputeMerkleRoot(leaf, [][32]byte{s0, s1}, 2, c.index)
		if got != c.want {
			t.Fatalf("index=%d: got %x, want %x", c.index, got, c.want)
		}
		if !VerifyMerkleBranch(leaf, [][32]byte{s0, s1}, 2, c.index, c.want) {
			t.Fatalf("index=%d: VerifyMerkleBranch should accept correct root", c.index)
		}
	}
}

func hashLevel(prevValue [32]byte, sibling [32]byte) [32]byte {
	return hash(prevValue, sibling)
}

func hashLevelSwap(prevValue [32]byte, sibling [32]byte) [32]byte {
	return hash(sibling, prevValue)
}

func TestVerifyMerkleBranchWrongRootFails(t *testing.T) {
	var leaf, s0, s1, wrongRoot [32]byte
	leaf[0], s0[0], s1[0], wrongRoot[0] = 0xAA, 0x01, 0x02, 0xFF

	if VerifyMerkleBranch(leaf, [][32]byte{s0, s1}, 2, 0, wrongRoot) {
		t.Fatal("VerifyMerkleBranch should reject a root that doesn't match")
	}
}

func TestVerifyMerkleBranchDepthZeroLeafIsRoot(t *testing.T) {
	var leaf [32]byte
	leaf[0] = 0xAA

	if !VerifyMerkleBranch(leaf, nil, 0, 0, leaf) {
		t.Fatal("depth 0: leaf itself is the root")
	}

	var wrongRoot [32]byte
	wrongRoot[0] = 0xBB
	if VerifyMerkleBranch(leaf, nil, 0, 0, wrongRoot) {
		t.Fatal("depth 0: mismatched root should fail")
	}
}

func TestVerifyMerkleBranchAtGindexMatchesPlainForm(t *testing.T) {
	var leaf, s0, s1 [32]byte
	leaf[0], s0[0], s1[0] = 0xAA, 0x01, 0x02
	branch := [][32]byte{s0, s1}

	// gindex=6 (0b110): depth=floorlog2(6)=2, subtree index=6-4=2 (0b10).
	const gindex = 6
	root := ComputeMerkleRoot(leaf, branch, 2, 2)

	if !VerifyMerkleBranchAtGindex(leaf, branch, gindex, root) {
		t.Fatal("VerifyMerkleBranchAtGindex should match VerifyMerkleBranch with derived depth/index")
	}

	var wrongRoot [32]byte
	wrongRoot[0] = 0xFF
	if VerifyMerkleBranchAtGindex(leaf, branch, gindex, wrongRoot) {
		t.Fatal("VerifyMerkleBranchAtGindex should reject a mismatched root")
	}
}

func TestFloorLog2(t *testing.T) {
	cases := []struct {
		x    uint64
		want uint64
	}{
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 2},
		{6, 2},
		{8, 3},
	}
	for _, c := range cases {
		if got := floorLog2(c.x); got != c.want {
			t.Fatalf("floorLog2(%d): got %d, want %d", c.x, got, c.want)
		}
	}
}

func TestSubtreeIndex(t *testing.T) {
	cases := []struct {
		gindex, depth uint64
		want          uint64
	}{
		{1, 0, 0},
		{2, 1, 0},
		{3, 1, 1},
		{6, 2, 2},
		{7, 2, 3},
	}
	for _, c := range cases {
		if got := subtreeIndex(c.gindex, c.depth); got != c.want {
			t.Fatalf("subtreeIndex(%d,%d): got %d, want %d", c.gindex, c.depth, got, c.want)
		}
	}
}
