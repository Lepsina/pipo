package consensus

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/Lepsina/pipo/internal/ssz"
)

func TestSyncCommitteeBitsMarshalUnmarshalRoundTrip(t *testing.T) {
	var b SyncCommitteeBits
	b[0], b[63] = 0x01, 0x02
	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got SyncCommitteeBits
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Fatalf("got %x, want %x", got, b)
	}
}

func TestSyncCommitteeBitsFixedSize(t *testing.T) {
	size, ok := SyncCommitteeBits{}.FixedSize()
	if !ok || size != 64 {
		t.Fatalf("got (%d,%v), want (64,true)", size, ok)
	}
}

func TestSyncCommitteePubkeysMarshalUnmarshalRoundTrip(t *testing.T) {
	var p SyncCommitteePubkeys
	p[0][0] = 0xAA
	p[511][0] = 0xBB

	data, err := p.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got SyncCommitteePubkeys
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != p {
		t.Fatal("round trip did not preserve value")
	}
}

func TestSyncCommitteePubkeysFixedSize(t *testing.T) {
	size, ok := SyncCommitteePubkeys{}.FixedSize()
	if !ok || size != 512*48 {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, 512*48)
	}
}

func TestSyncCommitteePubkeysUnmarshalWrongLengthErrors(t *testing.T) {
	var p SyncCommitteePubkeys
	if err := p.Unmarshal([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for wrong-length input")
	}
}

func TestCurrentSyncCommitteeBranchMarshalUnmarshalRoundTrip(t *testing.T) {
	var b CurrentSyncCommitteeBranch
	for i := range b {
		b[i][0] = byte(i + 1)
	}

	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got CurrentSyncCommitteeBranch
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Fatal("round trip did not preserve value")
	}
}

func TestCurrentSyncCommitteeBranchFixedSize(t *testing.T) {
	size, ok := CurrentSyncCommitteeBranch{}.FixedSize()
	if !ok || size != 6*32 {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, 6*32)
	}
}

func TestNextSyncCommitteeBranchMarshalUnmarshalRoundTrip(t *testing.T) {
	var b NextSyncCommitteeBranch
	for i := range b {
		b[i][0] = byte(i + 1)
	}

	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got NextSyncCommitteeBranch
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Fatal("round trip did not preserve value")
	}
}

func TestNextSyncCommitteeBranchFixedSize(t *testing.T) {
	size, ok := NextSyncCommitteeBranch{}.FixedSize()
	if !ok || size != 6*32 {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, 6*32)
	}
}

func TestFinalityBranchMarshalUnmarshalRoundTrip(t *testing.T) {
	var b FinalityBranch
	for i := range b {
		b[i][0] = byte(i + 1)
	}

	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got FinalityBranch
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Fatal("round trip did not preserve value")
	}
}

func TestFinalityBranchFixedSize(t *testing.T) {
	size, ok := FinalityBranch{}.FixedSize()
	if !ok || size != 7*32 {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, 7*32)
	}
}

func TestExecutionBranchMarshalUnmarshalRoundTrip(t *testing.T) {
	var b ExecutionBranch
	for i := range b {
		b[i][0] = byte(i + 1)
	}

	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got ExecutionBranch
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Fatal("round trip did not preserve value")
	}
}

func TestExecutionBranchFixedSize(t *testing.T) {
	size, ok := ExecutionBranch{}.FixedSize()
	if !ok || size != 4*32 {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, 4*32)
	}
}

func TestBeaconBlockHeaderMarshalFieldOrder(t *testing.T) {
	h := BeaconBlockHeader{
		Slot:          Slot(1),
		ProposerIndex: ValidatorIndex(2),
		ParentRoot:    Root{0x03},
		StateRoot:     Root{0x04},
		BodyRoot:      Root{0x05},
	}

	got, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var want []byte
	for _, m := range []ssz.Marshaler{h.Slot, h.ProposerIndex, h.ParentRoot, h.StateRoot, h.BodyRoot} {
		b, _ := m.Marshal()
		want = append(want, b...)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestBeaconBlockHeaderMarshalUnmarshalRoundTrip(t *testing.T) {
	h := BeaconBlockHeader{
		Slot:          Slot(1),
		ProposerIndex: ValidatorIndex(2),
		ParentRoot:    Root{0x03},
		StateRoot:     Root{0x04},
		BodyRoot:      Root{0x05},
	}

	data, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got BeaconBlockHeader
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != h {
		t.Fatalf("got %+v, want %+v", got, h)
	}
}

func TestBeaconBlockHeaderFixedSize(t *testing.T) {
	size, ok := BeaconBlockHeader{}.FixedSize()
	want := uint64(8 + 8 + 32 + 32 + 32)
	if !ok || size != want {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, want)
	}
}

func TestBeaconBlockHeaderUnmarshalTruncatedErrors(t *testing.T) {
	var h BeaconBlockHeader
	if err := h.Unmarshal(make([]byte, 50)); err == nil {
		t.Fatal("expected error for truncated input")
	}
}

func TestBeaconBlockHeaderUnmarshalTooLongErrors(t *testing.T) {
	var h BeaconBlockHeader
	if err := h.Unmarshal(make([]byte, 200)); err == nil {
		t.Fatal("expected error for oversized input")
	}
}

func TestSyncAggregateMarshalUnmarshalRoundTrip(t *testing.T) {
	a := SyncAggregate{
		SyncCommitteeBits:      SyncCommitteeBits{0x01},
		SyncCommitteeSignature: BLSSignature{0x02},
	}

	data, err := a.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got SyncAggregate
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != a {
		t.Fatalf("got %+v, want %+v", got, a)
	}
}

func TestSyncAggregateFixedSize(t *testing.T) {
	size, ok := SyncAggregate{}.FixedSize()
	want := uint64(64 + 96)
	if !ok || size != want {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, want)
	}
}

func TestSyncCommitteeMarshalUnmarshalRoundTrip(t *testing.T) {
	var c SyncCommittee
	c.Pubkeys[0][0] = 0xAA
	c.AggregatePubkey = BLSPubkey{0xBB}

	data, err := c.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got SyncCommittee
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != c {
		t.Fatal("round trip did not preserve value")
	}
}

func TestSyncCommitteeFixedSize(t *testing.T) {
	size, ok := SyncCommittee{}.FixedSize()
	want := uint64(512*48 + 48)
	if !ok || size != want {
		t.Fatalf("got (%d,%v), want (%d,true)", size, ok, want)
	}
}

func TestExtraDataMarshalUnmarshalRoundTrip(t *testing.T) {
	e := ExtraData{1, 2, 3}

	data, err := e.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got ExtraData
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(e) {
		t.Fatalf("got %x, want %x", got, e)
	}
}

func TestExtraDataMarshalEmpty(t *testing.T) {
	var e ExtraData
	data, err := e.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("got len %d, want 0", len(data))
	}
}

func TestExtraDataFixedSize(t *testing.T) {
	size, ok := ExtraData{}.FixedSize()
	if ok || size != 0 {
		t.Fatalf("got (%d,%v), want (0,false) -- ExtraData is variable-size", size, ok)
	}
}

func TestExtraDataUnmarshalAtMaxLengthSucceeds(t *testing.T) {
	var e ExtraData
	if err := e.Unmarshal(make([]byte, 32)); err != nil {
		t.Fatalf("unexpected error at max length: %v", err)
	}
}

func TestExtraDataUnmarshalTooLongErrors(t *testing.T) {
	var e ExtraData
	if err := e.Unmarshal(make([]byte, 33)); err == nil {
		t.Fatal("expected error: ExtraData exceeds MAX_EXTRA_DATA_BYTES=32")
	}
}

func TestExtraDataMarshalDoesNotAliasCaller(t *testing.T) {
	e := ExtraData{0xAA}
	data, _ := e.Marshal()
	data[0] = 0xFF
	if e[0] != 0xAA {
		t.Fatal("mutating the returned slice must not mutate the source ExtraData")
	}
}

func sampleExecutionPayloadHeader(extra []byte) ExecutionPayloadHeader {
	return ExecutionPayloadHeader{
		ParentHash:       Hash32{0x01},
		FeeRecipient:     ExecutionAddress{0x02},
		StateRoot:        Bytes32{0x03},
		ReceiptsRoot:     Bytes32{0x04},
		LogsBloom:        LogsBloom{0x05},
		PrevRandao:       Bytes32{0x06},
		BlockNumber:      Uint64(1),
		GasLimit:         Uint64(2),
		GasUsed:          Uint64(3),
		Timestamp:        Uint64(4),
		ExtraData:        ExtraData(extra),
		BaseFeePerGas:    Uint256{0x07},
		BlockHash:        Hash32{0x08},
		TransactionsRoot: Root{0x09},
		WithdrawalsRoot:  Root{0x0A},
		BlobGasUsed:      Uint64(5),
		ExcessBlobGas:    Uint64(6),
	}
}

// fixedPartLen = sum of every field's fixed size, with ExtraData (the one
// variable field) counted as a 4-byte offset instead of its real size.
const executionPayloadHeaderFixedPartLen = 32 + 20 + 32 + 32 + 256 + 32 + 8 + 8 + 8 + 8 + 4 + 32 + 32 + 32 + 32 + 8 + 8

// extraDataFixedPartOffset = sum of fixed sizes of every field before
// ExtraData (ParentHash..Timestamp).
const extraDataFixedPartOffset = 32 + 20 + 32 + 32 + 256 + 32 + 8 + 8 + 8 + 8

func TestExecutionPayloadHeaderFixedSizeIsVariable(t *testing.T) {
	size, ok := ExecutionPayloadHeader{}.FixedSize()
	if ok {
		t.Fatalf("ExecutionPayloadHeader embeds ExtraData, should be variable-size; got size=%d", size)
	}
}

func TestExecutionPayloadHeaderMarshalOffsetTable(t *testing.T) {
	h := sampleExecutionPayloadHeader([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE})

	data, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantLen := executionPayloadHeaderFixedPartLen + len(h.ExtraData)
	if len(data) != wantLen {
		t.Fatalf("got len %d, want %d", len(data), wantLen)
	}

	gotOffset := binary.LittleEndian.Uint32(data[extraDataFixedPartOffset : extraDataFixedPartOffset+4])
	if gotOffset != uint32(executionPayloadHeaderFixedPartLen) {
		t.Fatalf("got offset %d, want %d", gotOffset, executionPayloadHeaderFixedPartLen)
	}

	gotExtra := data[executionPayloadHeaderFixedPartLen:]
	if string(gotExtra) != string(h.ExtraData) {
		t.Fatalf("got extra_data %x, want %x", gotExtra, h.ExtraData)
	}
}

func TestExecutionPayloadHeaderMarshalUnmarshalRoundTrip(t *testing.T) {
	h := sampleExecutionPayloadHeader([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE})

	data, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got ExecutionPayloadHeader
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, h) {
		t.Fatalf("got %+v, want %+v", got, h)
	}
}

func TestExecutionPayloadHeaderMarshalUnmarshalRoundTripEmptyExtraData(t *testing.T) {
	h := sampleExecutionPayloadHeader(nil)

	data, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != executionPayloadHeaderFixedPartLen {
		t.Fatalf("got len %d, want %d (empty extra_data)", len(data), executionPayloadHeaderFixedPartLen)
	}

	var got ExecutionPayloadHeader
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.ExtraData) != 0 {
		t.Fatalf("got ExtraData len %d, want 0", len(got.ExtraData))
	}
}

func TestExecutionPayloadHeaderUnmarshalTruncatedFixedPartErrors(t *testing.T) {
	var h ExecutionPayloadHeader
	if err := h.Unmarshal(make([]byte, 10)); err == nil {
		t.Fatal("expected error for truncated fixed part")
	}
}

func TestExecutionPayloadHeaderUnmarshalBadOffsetErrors(t *testing.T) {
	h := sampleExecutionPayloadHeader([]byte{0xAA})
	data, _ := h.Marshal()

	// Corrupt the offset to point past the end of the buffer.
	binary.LittleEndian.PutUint32(data[extraDataFixedPartOffset:extraDataFixedPartOffset+4], uint32(len(data)+100))

	var got ExecutionPayloadHeader
	if err := got.Unmarshal(data); err == nil {
		t.Fatal("expected error for out-of-range offset")
	}
}

func TestLightClientHeaderFixedSizeIsVariable(t *testing.T) {
	if _, ok := (LightClientHeader{}).FixedSize(); ok {
		t.Fatal("LightClientHeader embeds ExecutionPayloadHeader, should be variable-size")
	}
}

func TestLightClientHeaderMarshalUnmarshalRoundTrip(t *testing.T) {
	h := sampleLightClientHeader(0x01)
	h.Execution.ExtraData = ExtraData{0xAA, 0xBB, 0xCC}

	data, err := h.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got LightClientHeader
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, h) {
		t.Fatalf("got %+v, want %+v", got, h)
	}
}

func TestLightClientBootstrapMarshalUnmarshalRoundTrip(t *testing.T) {
	h := sampleLightClientHeader(0x01)
	h.Execution.ExtraData = ExtraData{0xAA}

	b := LightClientBootstrap{
		Header:                     h,
		CurrentSyncCommittee:       SyncCommittee{AggregatePubkey: BLSPubkey{0x02}},
		CurrentSyncCommitteeBranch: CurrentSyncCommitteeBranch{{0x03}},
	}

	data, err := b.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got LightClientBootstrap
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, b) {
		t.Fatalf("got %+v, want %+v", got, b)
	}
}

func TestLightClientUpdateMarshalUnmarshalRoundTrip(t *testing.T) {
	attested := sampleLightClientHeader(0x01)
	attested.Execution.ExtraData = ExtraData{0xAA, 0xAA}
	finalized := sampleLightClientHeader(0x08)
	finalized.Execution.ExtraData = ExtraData{0xBB, 0xBB, 0xBB}

	u := LightClientUpdate{
		AttestedHeader:          attested,
		NextSyncCommittee:       SyncCommittee{AggregatePubkey: BLSPubkey{0x02}},
		NextSyncCommitteeBranch: NextSyncCommitteeBranch{{0x03}},
		FinalizedHeader:         finalized,
		FinalityBranch:          FinalityBranch{{0x09}},
		SyncAggregate:           SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:           Slot(99),
	}

	data, err := u.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got LightClientUpdate
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, u) {
		t.Fatalf("got %+v, want %+v", got, u)
	}
}

func TestLightClientFinalityUpdateMarshalUnmarshalRoundTrip(t *testing.T) {
	attested := sampleLightClientHeader(0x01)
	attested.Execution.ExtraData = ExtraData{0xAA}
	finalized := sampleLightClientHeader(0x08)
	finalized.Execution.ExtraData = ExtraData{0xBB, 0xBB}

	u := LightClientFinalityUpdate{
		AttestedHeader:  attested,
		FinalizedHeader: finalized,
		FinalityBranch:  FinalityBranch{{0x09}},
		SyncAggregate:   SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:   Slot(99),
	}

	data, err := u.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got LightClientFinalityUpdate
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, u) {
		t.Fatalf("got %+v, want %+v", got, u)
	}
}

func TestLightClientOptimisticUpdateMarshalUnmarshalRoundTrip(t *testing.T) {
	attested := sampleLightClientHeader(0x01)
	attested.Execution.ExtraData = ExtraData{0xAA, 0xAA, 0xAA}

	u := LightClientOptimisticUpdate{
		AttestedHeader: attested,
		SyncAggregate:  SyncAggregate{SyncCommitteeSignature: BLSSignature{0x0F}},
		SignatureSlot:  Slot(99),
	}

	data, err := u.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got LightClientOptimisticUpdate
	if err := got.Unmarshal(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, u) {
		t.Fatalf("got %+v, want %+v", got, u)
	}
}
