package bls

import (
	"errors"
	"math/big"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

const DST = "BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_"

var (
	ErrNilInput          = errors.New("bls: nil input")
	ErrEmptyInput        = errors.New("bls: empty input")
	ErrZeroKey           = errors.New("bls: zero secret key")
	ErrLenMismatch       = errors.New("bls: pubkeys and messages length mismatch")
	ErrTrailingBytes     = errors.New("bls: trailing bytes after encoded point")
	ErrNonCanonicalInfin = errors.New("bls: point at infinity encoded with a non-infinity header")
	ErrInfinityPubkey    = errors.New("bls: infinity point is not a valid public key")
	ErrInvalidLength     = errors.New("bls: encoded point has wrong length for its compressed form")
)

const infinityFlag = 0x40

type PublicKey struct{ p bls12381.G1Affine }
type Signature struct{ p bls12381.G2Affine }
type SecretKey struct{ s fr.Element }

func PublicKeyFromBytes(b []byte) (*PublicKey, error) {
	if len(b) != bls12381.SizeOfG1AffineCompressed {
		return nil, ErrInvalidLength
	}
	var pk bls12381.G1Affine
	n, err := pk.SetBytes(b)
	if err != nil {
		return nil, err
	}
	if n != len(b) {
		return nil, ErrTrailingBytes
	}
	if pk.IsInfinity() && b[0]&infinityFlag == 0 {
		return nil, ErrNonCanonicalInfin
	}
	return &PublicKey{p: pk}, nil
}

func SignatureFromBytes(b []byte) (*Signature, error) {
	if len(b) != bls12381.SizeOfG2AffineCompressed {
		return nil, ErrInvalidLength
	}
	var sig bls12381.G2Affine
	n, err := sig.SetBytes(b)
	if err != nil {
		return nil, err
	}
	if n != len(b) {
		return nil, ErrTrailingBytes
	}
	if sig.IsInfinity() && b[0]&infinityFlag == 0 {
		return nil, ErrNonCanonicalInfin
	}
	return &Signature{p: sig}, nil
}

// SecretKeyFromBytes parses a 32-byte big-endian scalar in the range [1, r-1].
func SecretKeyFromBytes(b []byte) (*SecretKey, error) {
	var s fr.Element
	if err := s.SetBytesCanonical(b); err != nil {
		return nil, err
	}
	if s.IsZero() {
		return nil, ErrZeroKey
	}
	return &SecretKey{s: s}, nil
}

func (pk *PublicKey) Bytes() [bls12381.SizeOfG1AffineCompressed]byte  { return pk.p.Bytes() }
func (sig *Signature) Bytes() [bls12381.SizeOfG2AffineCompressed]byte { return sig.p.Bytes() }

// KeyValidate reports whether b encodes a valid public key:
// a well-formed point in the G1 subgroup that is not the identity.
func KeyValidate(b []byte) (bool, error) {
	pk, err := PublicKeyFromBytes(b)
	if err != nil {
		return false, err
	}
	if pk.p.IsInfinity() {
		return false, ErrInfinityPubkey
	}
	return true, nil
}

func (sk *SecretKey) PublicKey() *PublicKey {
	_, _, g1Gen, _ := bls12381.Generators()
	var pk bls12381.G1Affine
	pk.ScalarMultiplication(&g1Gen, sk.s.BigInt(new(big.Int)))
	return &PublicKey{p: pk}
}

func (sk *SecretKey) Sign(msg []byte) (*Signature, error) {
	Hm, err := bls12381.HashToG2(msg, []byte(DST))
	if err != nil {
		return nil, err
	}
	var sig bls12381.G2Affine
	sig.ScalarMultiplication(&Hm, sk.s.BigInt(new(big.Int)))
	return &Signature{p: sig}, nil
}

func Verify(pk *PublicKey, msg []byte, sig *Signature) (bool, error) {
	if pk == nil || sig == nil || pk.p.IsInfinity() {
		return false, nil
	}

	Hm, err := bls12381.HashToG2(msg, []byte(DST))
	if err != nil {
		return false, err
	}

	_, _, g1Gen, _ := bls12381.Generators()
	var negG1 bls12381.G1Affine
	negG1.Neg(&g1Gen)

	// e(pk, H(m)) * e(-G1, sig) == 1
	return bls12381.PairingCheck(
		[]bls12381.G1Affine{pk.p, negG1},
		[]bls12381.G2Affine{Hm, sig.p},
	)
}

func Aggregate(sigs []*Signature) (*Signature, error) {
	if len(sigs) == 0 {
		return nil, ErrEmptyInput
	}
	var agg bls12381.G2Jac
	for _, s := range sigs {
		if s == nil {
			return nil, ErrNilInput
		}
		var j bls12381.G2Jac
		j.FromAffine(&s.p)
		agg.AddAssign(&j)
	}
	var out bls12381.G2Affine
	out.FromJacobian(&agg)
	return &Signature{p: out}, nil
}

func AggregatePublicKeys(pks []*PublicKey) (*PublicKey, error) {
	if len(pks) == 0 {
		return nil, ErrEmptyInput
	}
	var agg bls12381.G1Jac
	for _, pk := range pks {
		if pk == nil {
			return nil, ErrNilInput
		}
		var j bls12381.G1Jac
		j.FromAffine(&pk.p)
		agg.AddAssign(&j)
	}
	var out bls12381.G1Affine
	out.FromJacobian(&agg)
	return &PublicKey{p: out}, nil
}

// AggregateVerify checks an aggregate signature over per-key messages
// in the proof-of-possession scheme.
func AggregateVerify(pks []*PublicKey, msgs [][]byte, sig *Signature) (bool, error) {
	if len(pks) != len(msgs) {
		return false, ErrLenMismatch
	}
	if len(pks) == 0 || sig == nil {
		return false, nil
	}

	g1s := make([]bls12381.G1Affine, 0, len(pks)+1)
	g2s := make([]bls12381.G2Affine, 0, len(pks)+1)
	for i, pk := range pks {
		if pk == nil || pk.p.IsInfinity() {
			return false, nil
		}
		Hm, err := bls12381.HashToG2(msgs[i], []byte(DST))
		if err != nil {
			return false, err
		}
		g1s = append(g1s, pk.p)
		g2s = append(g2s, Hm)
	}

	_, _, g1Gen, _ := bls12381.Generators()
	var negG1 bls12381.G1Affine
	negG1.Neg(&g1Gen)
	g1s = append(g1s, negG1)
	g2s = append(g2s, sig.p)

	return bls12381.PairingCheck(g1s, g2s)
}

func FastAggregateVerify(pks []*PublicKey, msg []byte, sig *Signature) (bool, error) {
	if len(pks) == 0 {
		return false, nil
	}
	for _, pk := range pks {
		if pk == nil || pk.p.IsInfinity() {
			return false, nil
		}
	}
	aggregatedPK, err := AggregatePublicKeys(pks)
	if err != nil {
		return false, err
	}
	return Verify(aggregatedPK, msg, sig)
}

// EthFastAggregateVerify is the Ethereum consensus-spec variant: an empty
// pubkey set is valid if the signature is the point at infinity.
func EthFastAggregateVerify(pks []*PublicKey, msg []byte, sig *Signature) (bool, error) {
	if sig == nil {
		return false, nil
	}
	if len(pks) == 0 && sig.p.IsInfinity() {
		return true, nil
	}
	return FastAggregateVerify(pks, msg, sig)
}
