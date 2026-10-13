package consensus

import "github.com/Lepsina/pipo/internal/ssz"

func (b BeaconBlockHeader) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(b)
}

func (a SyncAggregate) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(a)
}

func (c SyncCommittee) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(c)
}

func (h ExecutionPayloadHeader) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(h)
}

func (h LightClientHeader) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(h)
}

func (b LightClientBootstrap) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(b)
}

func (u LightClientUpdate) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(u)
}

func (u LightClientFinalityUpdate) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(u)
}

func (u LightClientOptimisticUpdate) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootContainer(u)
}

func (b SyncCommitteeBits) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootByteVector(b[:])
}

func (p SyncCommitteePubkeys) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootVector(p[:])
}

func (b CurrentSyncCommitteeBranch) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootVector(b[:])
}

func (b FinalityBranch) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootVector(b[:])
}

func (b NextSyncCommitteeBranch) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootVector(b[:])
}

func (b ExecutionBranch) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootVector(b[:])
}

func (e ExtraData) HashTreeRoot() [32]byte {
	return ssz.HashTreeRootByteList(e, maxExtraDataBytes)
}
