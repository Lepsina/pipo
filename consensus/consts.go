package consensus

// Consensus constants that fix the sizes of this package's types.
const (
	// TODO: add a testnet preset
	// mainnet preset
	syncCommitteeSize = 512

	// SYNC_COMMITTEE_SIZE bits, packed 8 to a byte.
	syncCommitteeBitsSize = syncCommitteeSize / 8

	maxExtraDataBytes = 32
)
