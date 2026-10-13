// Package consensus defines the consensus-layer containers the light client
// works with. Each type declares its SSZ shape here; the SSZ rules
// themselves live in internal/ssz.
package consensus

import "github.com/Lepsina/pipo/internal/ssz"

type Uint64 = ssz.Uint64

type Uint256 = ssz.Uint256

type Bytes32 = ssz.Bytes32

type Root = Bytes32

type Hash32 = Bytes32

type Slot = Uint64

type ValidatorIndex = Uint64

type BLSPubkey = ssz.Bytes48

type BLSSignature = ssz.Bytes96

type ExecutionAddress = ssz.Bytes20

type LogsBloom = ssz.Bytes256

type BeaconBlockHeader struct {
	Slot          Slot
	ProposerIndex ValidatorIndex
	ParentRoot    Root
	StateRoot     Root
	BodyRoot      Root
}

type SyncAggregate struct {
	SyncCommitteeBits      SyncCommitteeBits
	SyncCommitteeSignature BLSSignature
}

type SyncCommittee struct {
	Pubkeys         SyncCommitteePubkeys
	AggregatePubkey BLSPubkey
}

// SyncCommitteeBits is Bitvector[SYNC_COMMITTEE_SIZE].
type SyncCommitteeBits [syncCommitteeBitsSize]byte

// SyncCommitteePubkeys is Vector[BLSPubkey, SYNC_COMMITTEE_SIZE].
type SyncCommitteePubkeys [syncCommitteeSize]BLSPubkey

// LENGTH = floorlog2(CURRENT_SYNC_COMMITTEE_GINDEX_ELECTRA) = floorlog2(86) = 6.
type CurrentSyncCommitteeBranch [6]Bytes32

// LENGTH = floorlog2(FINALIZED_ROOT_GINDEX_ELECTRA) = floorlog2(169) = 7.
type FinalityBranch [7]Bytes32

// LENGTH = floorlog2(NEXT_SYNC_COMMITTEE_GINDEX_ELECTRA) = floorlog2(87) = 6.
type NextSyncCommitteeBranch [6]Bytes32

// LENGTH = floorlog2(EXECUTION_PAYLOAD_GINDEX) = floorlog2(25) = 4.
type ExecutionBranch [4]Bytes32

// ExtraData is ByteList[MAX_EXTRA_DATA_BYTES].
type ExtraData []byte

// ExecutionPayloadHeader, Deneb shape (current through Electra/Fulu — unchanged since Deneb).
type ExecutionPayloadHeader struct {
	ParentHash       Hash32
	FeeRecipient     ExecutionAddress
	StateRoot        Bytes32
	ReceiptsRoot     Bytes32
	LogsBloom        LogsBloom
	PrevRandao       Bytes32
	BlockNumber      Uint64
	GasLimit         Uint64
	GasUsed          Uint64
	Timestamp        Uint64
	ExtraData        ExtraData
	BaseFeePerGas    Uint256
	BlockHash        Hash32
	TransactionsRoot Root
	WithdrawalsRoot  Root
	BlobGasUsed      Uint64
	ExcessBlobGas    Uint64
}

type LightClientHeader struct {
	Beacon          BeaconBlockHeader
	Execution       ExecutionPayloadHeader
	ExecutionBranch ExecutionBranch
}

type LightClientBootstrap struct {
	Header                     LightClientHeader
	CurrentSyncCommittee       SyncCommittee
	CurrentSyncCommitteeBranch CurrentSyncCommitteeBranch
}

type LightClientUpdate struct {
	AttestedHeader          LightClientHeader
	NextSyncCommittee       SyncCommittee
	NextSyncCommitteeBranch NextSyncCommitteeBranch
	FinalizedHeader         LightClientHeader
	FinalityBranch          FinalityBranch
	SyncAggregate           SyncAggregate
	SignatureSlot           Slot
}

type LightClientFinalityUpdate struct {
	AttestedHeader  LightClientHeader
	FinalizedHeader LightClientHeader
	FinalityBranch  FinalityBranch
	SyncAggregate   SyncAggregate
	SignatureSlot   Slot
}

type LightClientOptimisticUpdate struct {
	AttestedHeader LightClientHeader
	SyncAggregate  SyncAggregate
	SignatureSlot  Slot
}
