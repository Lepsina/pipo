//go:build blst

// Differential fuzzing against github.com/supranational/blst
package bls

import (
	"bytes"
	"crypto/sha256"
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	blst "github.com/supranational/blst/bindings/go"
)

func FuzzDiffDecode(f *testing.F) {
	seedFromField(f, "deserialization_G1", "pubkey")
	seedFromField(f, "deserialization_G2", "signature")

	f.Fuzz(func(t *testing.T, data []byte) {
		diffDecodeG1(t, data)
		diffDecodeG2(t, data)
	})
}

func diffDecodeG1(t *testing.T, data []byte) {
	t.Helper()

	gpk, gErr := PublicKeyFromBytes(data)
	gOK := gErr == nil

	if len(data) != bls12381.SizeOfG1AffineCompressed {
		if gOK {
			t.Fatalf("G1: accepted a non-compressed length %d: %x", len(data), data)
		}
		return
	}

	bp := new(blst.P1Affine).Uncompress(data)
	bOK := bp != nil && bp.InG1()

	if gOK != bOK {
		t.Fatalf("G1 decode disagreement: gnark=%v blst=%v data=%x", gOK, bOK, data)
	}
	if !gOK {
		return
	}
	genc := gpk.Bytes()
	if benc := bp.Compress(); !bytes.Equal(genc[:], benc) {
		t.Fatalf("G1 decode point mismatch: gnark=%x blst=%x", genc, benc)
	}
}

func diffDecodeG2(t *testing.T, data []byte) {
	t.Helper()

	gsig, gErr := SignatureFromBytes(data)
	gOK := gErr == nil

	if len(data) != bls12381.SizeOfG2AffineCompressed {
		if gOK {
			t.Fatalf("G2: accepted a non-compressed length %d: %x", len(data), data)
		}
		return
	}

	bp := new(blst.P2Affine).Uncompress(data)
	bOK := bp != nil && bp.InG2()

	if gOK != bOK {
		t.Fatalf("G2 decode disagreement: gnark=%v blst=%v data=%x", gOK, bOK, data)
	}
	if !gOK {
		return
	}
	genc := gsig.Bytes()
	if benc := bp.Compress(); !bytes.Equal(genc[:], benc) {
		t.Fatalf("G2 decode point mismatch: gnark=%x blst=%x", genc, benc)
	}
}

func deriveSK(t *testing.T, seed []byte, tag byte) (*SecretKey, *blst.SecretKey, bool) {
	t.Helper()

	h := sha256.Sum256(append([]byte{tag}, seed...))
	gsk, gErr := SecretKeyFromBytes(h[:])
	bsk := new(blst.SecretKey).Deserialize(h[:])

	if (gErr == nil) != (bsk != nil) {
		t.Fatalf("secret key validity disagreement for %x: gnark=%v blst=%v", h, gErr == nil, bsk != nil)
	}
	return gsk, bsk, gErr == nil
}

func FuzzDiffVerify(f *testing.F) {
	f.Add([]byte("seed"), []byte("message"), false, byte(0))
	f.Add([]byte("seed"), []byte("message"), true, byte(3))

	f.Fuzz(func(t *testing.T, seed, msg []byte, corrupt bool, corruptIdx byte) {
		const n = 3
		gsks := make([]*SecretKey, 0, n)
		bsks := make([]*blst.SecretKey, 0, n)
		for i := byte(0); i < n; i++ {
			gsk, bsk, ok := deriveSK(t, seed, i)
			if !ok {
				return // both libraries reject this seed family; nothing to verify
			}
			gsks = append(gsks, gsk)
			bsks = append(bsks, bsk)
		}
		_ = bsks

		pks := make([]*PublicKey, 0, n)
		sigs := make([]*Signature, 0, n)
		for _, gsk := range gsks {
			pks = append(pks, gsk.PublicKey())
			sig, err := gsk.Sign(msg)
			if err != nil {
				t.Fatalf("sign: %v", err)
			}
			sigs = append(sigs, sig)
		}

		aggSig, err := Aggregate(sigs)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		sigBytes := aggSig.Bytes()
		if corrupt {
			sigBytes[int(corruptIdx)%len(sigBytes)] ^= 0xff
		}

		gSig, gErr := SignatureFromBytes(sigBytes[:])
		bSig := new(blst.P2Affine).Uncompress(sigBytes[:])
		bSigOK := bSig != nil && bSig.InG2()
		if (gErr == nil) != bSigOK {
			t.Fatalf("(possibly corrupted) signature decode disagreement: gnark=%v blst=%v corrupt=%v",
				gErr == nil, bSigOK, corrupt)
		}
		if gErr != nil {
			return
		}

		bpks := make([]*blst.P1Affine, 0, n)
		for _, pk := range pks {
			enc := pk.Bytes()
			bpk := new(blst.P1Affine).Uncompress(enc[:])
			if bpk == nil || !bpk.InG1() {
				t.Fatalf("gnark-valid pubkey rejected by blst: %x", enc)
			}
			bpks = append(bpks, bpk)
		}

		gOK, err := FastAggregateVerify(pks, msg, gSig)
		if err != nil {
			t.Fatalf("gnark FastAggregateVerify: %v", err)
		}
		bOK := bSig.FastAggregateVerify(true, bpks, msg, []byte(DST))
		if gOK != bOK {
			t.Fatalf("FastAggregateVerify disagreement: gnark=%v blst=%v corrupt=%v", gOK, bOK, corrupt)
		}

		gVOK, err := Verify(pks[0], msg, sigs[0])
		if err != nil {
			t.Fatalf("gnark Verify: %v", err)
		}
		pkEnc := pks[0].Bytes()
		sigEnc := sigs[0].Bytes()
		bpk0 := new(blst.P1Affine).Uncompress(pkEnc[:])
		bsig0 := new(blst.P2Affine).Uncompress(sigEnc[:])
		bVOK := bsig0.Verify(true, bpk0, true, msg, []byte(DST))
		if gVOK != bVOK {
			t.Fatalf("Verify disagreement: gnark=%v blst=%v", gVOK, bVOK)
		}
	})
}
