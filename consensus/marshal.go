package consensus

import "github.com/Lepsina/pipo/internal/ssz"

// ---- Marshal ----

func (b BeaconBlockHeader) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(b)
}

func (a SyncAggregate) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(a)
}

func (c SyncCommittee) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(c)
}

func (h ExecutionPayloadHeader) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(h)
}

func (h LightClientHeader) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(h)
}

func (b LightClientBootstrap) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(b)
}

func (u LightClientUpdate) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(u)
}

func (u LightClientFinalityUpdate) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(u)
}

func (u LightClientOptimisticUpdate) Marshal() ([]byte, error) {
	return ssz.MarshalContainer(u)
}

func (b SyncCommitteeBits) Marshal() ([]byte, error) {
	return ssz.MarshalByteVector(b[:])
}

func (p SyncCommitteePubkeys) Marshal() ([]byte, error) {
	return ssz.MarshalVector(p[:])
}

func (b CurrentSyncCommitteeBranch) Marshal() ([]byte, error) {
	return ssz.MarshalVector(b[:])
}

func (b NextSyncCommitteeBranch) Marshal() ([]byte, error) {
	return ssz.MarshalVector(b[:])
}

func (b FinalityBranch) Marshal() ([]byte, error) {
	return ssz.MarshalVector(b[:])
}

func (b ExecutionBranch) Marshal() ([]byte, error) {
	return ssz.MarshalVector(b[:])
}

func (e ExtraData) Marshal() ([]byte, error) {
	return ssz.MarshalByteList(e)
}

// ---- Unmarshal ----

func (b *BeaconBlockHeader) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(b, data)
}

func (a *SyncAggregate) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(a, data)
}

func (c *SyncCommittee) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(c, data)
}

func (h *ExecutionPayloadHeader) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(h, data)
}

func (h *LightClientHeader) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(h, data)
}

func (b *LightClientBootstrap) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(b, data)
}

func (u *LightClientUpdate) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(u, data)
}

func (u *LightClientFinalityUpdate) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(u, data)
}

func (u *LightClientOptimisticUpdate) Unmarshal(data []byte) error {
	return ssz.UnmarshalContainer(u, data)
}

func (b *SyncCommitteeBits) Unmarshal(data []byte) error {
	return ssz.UnmarshalByteVector(b[:], data)
}

func (p *SyncCommitteePubkeys) Unmarshal(data []byte) error {
	elemSize, _ := BLSPubkey{}.FixedSize()
	return ssz.UnmarshalVector(p[:], data, elemSize)
}

func (b *CurrentSyncCommitteeBranch) Unmarshal(data []byte) error {
	elemSize, _ := Root{}.FixedSize()
	return ssz.UnmarshalVector(b[:], data, elemSize)
}

func (b *NextSyncCommitteeBranch) Unmarshal(data []byte) error {
	elemSize, _ := Root{}.FixedSize()
	return ssz.UnmarshalVector(b[:], data, elemSize)
}

func (b *FinalityBranch) Unmarshal(data []byte) error {
	elemSize, _ := Root{}.FixedSize()
	return ssz.UnmarshalVector(b[:], data, elemSize)
}

func (b *ExecutionBranch) Unmarshal(data []byte) error {
	elemSize, _ := Root{}.FixedSize()
	return ssz.UnmarshalVector(b[:], data, elemSize)
}

func (e *ExtraData) Unmarshal(data []byte) error {
	b, err := ssz.UnmarshalByteList(data, maxExtraDataBytes)
	if err != nil {
		return err
	}
	*e = b
	return nil
}

// ---- FixedSize ----

func (b BeaconBlockHeader) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(b)
}

func (a SyncAggregate) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(a)
}

func (c SyncCommittee) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(c)
}

func (h ExecutionPayloadHeader) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(h)
}

func (h LightClientHeader) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(h)
}

func (b LightClientBootstrap) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(b)
}

func (u LightClientUpdate) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(u)
}

func (u LightClientFinalityUpdate) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(u)
}

func (u LightClientOptimisticUpdate) FixedSize() (uint64, bool) {
	return ssz.FixedSizeContainer(u)
}

func (SyncCommitteeBits) FixedSize() (uint64, bool) { return syncCommitteeBitsSize, true }

func (p SyncCommitteePubkeys) FixedSize() (uint64, bool) {
	elemSize, _ := BLSPubkey{}.FixedSize()
	return uint64(len(p)) * elemSize, true
}

func (b CurrentSyncCommitteeBranch) FixedSize() (uint64, bool) {
	elemSize, _ := Root{}.FixedSize()
	return uint64(len(b)) * elemSize, true
}

func (b NextSyncCommitteeBranch) FixedSize() (uint64, bool) {
	elemSize, _ := Root{}.FixedSize()
	return uint64(len(b)) * elemSize, true
}

func (b FinalityBranch) FixedSize() (uint64, bool) {
	elemSize, _ := Root{}.FixedSize()
	return uint64(len(b)) * elemSize, true
}

func (b ExecutionBranch) FixedSize() (uint64, bool) {
	elemSize, _ := Root{}.FixedSize()
	return uint64(len(b)) * elemSize, true
}

func (ExtraData) FixedSize() (uint64, bool) { return 0, false }
