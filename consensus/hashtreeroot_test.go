package consensus

import (
	"encoding/binary"
	"testing"

	"github.com/Lepsina/pipo/internal/ssz"
)

// rootOfFieldRoots is the Merkle root over a container's field roots, which
// a container's hash_tree_root must equal. It goes through the vector rule
// rather than the container rule, so the expected value comes from a
// different code path than the one under test.
func rootOfFieldRoots(roots ...[32]byte) [32]byte {
	rs := make([]Root, len(roots))
	for i, r := range roots {
		rs[i] = r
	}
	return ssz.HashTreeRootVector(rs)
}

func TestSlotHashTreeRootIsPackedValueUnchanged(t *testing.T) {
	s := Slot(5)

	var want [32]byte
	binary.LittleEndian.PutUint64(want[:8], uint64(s))

	got := s.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestValidatorIndexHashTreeRootIsPackedValueUnchanged(t *testing.T) {
	v := ValidatorIndex(42)

	var want [32]byte
	binary.LittleEndian.PutUint64(want[:8], uint64(v))

	got := v.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBeaconBlockHeaderHashTreeRoot(t *testing.T) {
	h := BeaconBlockHeader{
		Slot:          5,
		ProposerIndex: 7,
		ParentRoot:    Root{0x01},
		StateRoot:     Root{0x02},
		BodyRoot:      Root{0x03},
	}

	chunks := [][32]byte{
		h.Slot.HashTreeRoot(),
		h.ProposerIndex.HashTreeRoot(),
		h.ParentRoot.HashTreeRoot(),
		h.StateRoot.HashTreeRoot(),
		h.BodyRoot.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := h.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestLightClientHeaderHashTreeRoot(t *testing.T) {
	h := LightClientHeader{
		Beacon: BeaconBlockHeader{
			Slot:          5,
			ProposerIndex: 7,
			ParentRoot:    Root{0x01},
			StateRoot:     Root{0x02},
			BodyRoot:      Root{0x03},
		},
		Execution: ExecutionPayloadHeader{
			ParentHash:   Hash32{0x11},
			FeeRecipient: ExecutionAddress{0x12},
		},
		ExecutionBranch: ExecutionBranch{{0x21}, {0x22}, {0x23}, {0x24}},
	}

	chunks := [][32]byte{
		h.Beacon.HashTreeRoot(),
		h.Execution.HashTreeRoot(),
		h.ExecutionBranch.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := h.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func sampleLightClientHeader(seed byte) LightClientHeader {
	return LightClientHeader{
		Beacon: BeaconBlockHeader{
			Slot:          Slot(seed),
			ProposerIndex: ValidatorIndex(seed),
			ParentRoot:    Root{seed},
			StateRoot:     Root{seed + 1},
			BodyRoot:      Root{seed + 2},
		},
		Execution: ExecutionPayloadHeader{
			ParentHash:   Hash32{seed},
			FeeRecipient: ExecutionAddress{seed},
		},
		ExecutionBranch: ExecutionBranch{{seed}, {seed + 1}, {seed + 2}, {seed + 3}},
	}
}

func TestLightClientBootstrapHashTreeRoot(t *testing.T) {
	b := LightClientBootstrap{
		Header:                     sampleLightClientHeader(0x01),
		CurrentSyncCommittee:       SyncCommittee{AggregatePubkey: BLSPubkey{0x02}},
		CurrentSyncCommitteeBranch: CurrentSyncCommitteeBranch{{0x03}, {0x04}, {0x05}, {0x06}, {0x07}, {0x08}},
	}

	chunks := [][32]byte{
		b.Header.HashTreeRoot(),
		b.CurrentSyncCommittee.HashTreeRoot(),
		b.CurrentSyncCommitteeBranch.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := b.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestLightClientUpdateHashTreeRoot(t *testing.T) {
	u := LightClientUpdate{
		AttestedHeader:          sampleLightClientHeader(0x01),
		NextSyncCommittee:       SyncCommittee{AggregatePubkey: BLSPubkey{0x02}},
		NextSyncCommitteeBranch: NextSyncCommitteeBranch{{0x03}, {0x04}, {0x05}, {0x06}, {0x07}, {0x08}},
		FinalizedHeader:         sampleLightClientHeader(0x08),
		FinalityBranch:          FinalityBranch{{0x09}, {0x0A}, {0x0B}, {0x0C}, {0x0D}, {0x0E}, {0x0F}},
		SyncAggregate:           SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:           Slot(99),
	}

	chunks := [][32]byte{
		u.AttestedHeader.HashTreeRoot(),
		u.NextSyncCommittee.HashTreeRoot(),
		u.NextSyncCommitteeBranch.HashTreeRoot(),
		u.FinalizedHeader.HashTreeRoot(),
		u.FinalityBranch.HashTreeRoot(),
		u.SyncAggregate.HashTreeRoot(),
		u.SignatureSlot.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := u.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestLightClientFinalityUpdateHashTreeRoot(t *testing.T) {
	u := LightClientFinalityUpdate{
		AttestedHeader:  sampleLightClientHeader(0x01),
		FinalizedHeader: sampleLightClientHeader(0x08),
		FinalityBranch:  FinalityBranch{{0x09}, {0x0A}, {0x0B}, {0x0C}, {0x0D}, {0x0E}, {0x0F}},
		SyncAggregate:   SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:   Slot(99),
	}

	chunks := [][32]byte{
		u.AttestedHeader.HashTreeRoot(),
		u.FinalizedHeader.HashTreeRoot(),
		u.FinalityBranch.HashTreeRoot(),
		u.SyncAggregate.HashTreeRoot(),
		u.SignatureSlot.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := u.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestLightClientOptimisticUpdateHashTreeRoot(t *testing.T) {
	u := LightClientOptimisticUpdate{
		AttestedHeader: sampleLightClientHeader(0x01),
		SyncAggregate:  SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:  Slot(99),
	}

	chunks := [][32]byte{
		u.AttestedHeader.HashTreeRoot(),
		u.SyncAggregate.HashTreeRoot(),
		u.SignatureSlot.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := u.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestExecutionPayloadHeaderHashTreeRoot(t *testing.T) {
	var bloom LogsBloom
	bloom[0] = 0x01

	h := ExecutionPayloadHeader{
		ParentHash:       Hash32{0x01},
		FeeRecipient:     ExecutionAddress{0x02},
		StateRoot:        Bytes32{0x03},
		ReceiptsRoot:     Bytes32{0x04},
		LogsBloom:        bloom,
		PrevRandao:       Bytes32{0x05},
		BlockNumber:      100,
		GasLimit:         200,
		GasUsed:          150,
		Timestamp:        12345,
		ExtraData:        ExtraData{0x06, 0x07},
		BaseFeePerGas:    Uint256{0x08},
		BlockHash:        Hash32{0x09},
		TransactionsRoot: Root{0x0A},
		WithdrawalsRoot:  Root{0x0B},
		BlobGasUsed:      300,
		ExcessBlobGas:    400,
	}

	chunks := [][32]byte{
		h.ParentHash.HashTreeRoot(),
		h.FeeRecipient.HashTreeRoot(),
		h.StateRoot.HashTreeRoot(),
		h.ReceiptsRoot.HashTreeRoot(),
		h.LogsBloom.HashTreeRoot(),
		h.PrevRandao.HashTreeRoot(),
		h.BlockNumber.HashTreeRoot(),
		h.GasLimit.HashTreeRoot(),
		h.GasUsed.HashTreeRoot(),
		h.Timestamp.HashTreeRoot(),
		h.ExtraData.HashTreeRoot(),
		h.BaseFeePerGas.HashTreeRoot(),
		h.BlockHash.HashTreeRoot(),
		h.TransactionsRoot.HashTreeRoot(),
		h.WithdrawalsRoot.HashTreeRoot(),
		h.BlobGasUsed.HashTreeRoot(),
		h.ExcessBlobGas.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := h.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestSyncAggregateHashTreeRoot(t *testing.T) {
	var bits SyncCommitteeBits
	bits[0] = 0xFF
	a := SyncAggregate{
		SyncCommitteeBits:      bits,
		SyncCommitteeSignature: BLSSignature{0x01},
	}

	chunks := [][32]byte{
		a.SyncCommitteeBits.HashTreeRoot(),
		a.SyncCommitteeSignature.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := a.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestSyncCommitteeHashTreeRoot(t *testing.T) {
	var pubkeys SyncCommitteePubkeys
	pubkeys[0][0] = 0xAB
	c := SyncCommittee{
		Pubkeys:         pubkeys,
		AggregatePubkey: BLSPubkey{0x01},
	}

	chunks := [][32]byte{
		c.Pubkeys.HashTreeRoot(),
		c.AggregatePubkey.HashTreeRoot(),
	}
	want := rootOfFieldRoots(chunks...)

	got := c.HashTreeRoot()

	if got != want {
		t.Fatalf("got %x, want %x", got, want)
	}
}
