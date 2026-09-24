package bls

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"gopkg.in/yaml.v3"
)

const vectorsDir = "../testdata/bls-12-381-tests"

type vector struct {
	Input  any `yaml:"input"`
	Output any `yaml:"output"`
}

type namedVector struct {
	name string
	v    vector
}

func loadVectors(dir, group string) ([]namedVector, error) {
	files, err := filepath.Glob(filepath.Join(dir, group, "*.yaml"))
	if err != nil {
		return nil, err
	}
	out := make([]namedVector, 0, len(files))
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var v vector
		if err := yaml.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		out = append(out, namedVector{name: strings.TrimSuffix(filepath.Base(f), ".yaml"), v: v})
	}
	return out, nil
}

func forEachVector(t *testing.T, group string, fn func(t *testing.T, v vector)) {
	t.Helper()
	vs, err := loadVectors(vectorsDir, group)
	if err != nil || len(vs) == 0 {
		t.Fatalf("no vectors found for %s: %v", group, err)
	}
	for _, nv := range vs {
		t.Run(nv.name, func(t *testing.T) { fn(t, nv.v) })
	}
}

func in(v vector) map[string]any { return v.Input.(map[string]any) }

func hexBytes(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}

func unhex(t *testing.T, x any) []byte {
	t.Helper()
	s, ok := x.(string)
	if !ok {
		t.Fatalf("expected hex string, got %T", x)
	}
	b, err := hexBytes(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func seedFromField(f *testing.F, group, field string) {
	f.Helper()
	vs, err := loadVectors(vectorsDir, group)
	if err != nil {
		f.Fatalf("load vectors for %s: %v", group, err)
	}
	for _, nv := range vs {
		s, ok := in(nv.v)[field].(string)
		if !ok {
			continue
		}
		data, err := hexBytes(s)
		if err != nil {
			continue
		}
		f.Add(data)
	}
}

func unhexList(t *testing.T, x any) [][]byte {
	t.Helper()
	l, ok := x.([]any)
	if !ok {
		t.Fatalf("expected list, got %T", x)
	}
	out := make([][]byte, len(l))
	for i := range l {
		out[i] = unhex(t, l[i])
	}
	return out
}

func parsePubkeys(t *testing.T, raw [][]byte) ([]*PublicKey, bool) {
	t.Helper()
	pks := make([]*PublicKey, len(raw))
	for i, b := range raw {
		pk, err := PublicKeyFromBytes(b)
		if err != nil {
			return nil, false
		}
		pks[i] = pk
	}
	return pks, true
}

func TestSign(t *testing.T) {
	forEachVector(t, "sign", func(t *testing.T, v vector) {
		sk, err := SecretKeyFromBytes(unhex(t, in(v)["privkey"]))
		if v.Output == nil {
			if err == nil {
				t.Fatal("expected error for invalid secret key")
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		sig, err := sk.Sign(unhex(t, in(v)["message"]))
		if err != nil {
			t.Fatal(err)
		}
		got := sig.Bytes()
		if want := unhex(t, v.Output); string(got[:]) != string(want) {
			t.Fatalf("signature mismatch:\n got  %x\n want %x", got, want)
		}
	})
}

func TestVerify(t *testing.T) {
	forEachVector(t, "verify", func(t *testing.T, v vector) {
		want := v.Output.(bool)
		pk, err := PublicKeyFromBytes(unhex(t, in(v)["pubkey"]))
		if err != nil {
			if want {
				t.Fatal(err)
			}
			return
		}
		sig, err := SignatureFromBytes(unhex(t, in(v)["signature"]))
		if err != nil {
			if want {
				t.Fatal(err)
			}
			return
		}
		got, err := Verify(pk, unhex(t, in(v)["message"]), sig)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestAggregate(t *testing.T) {
	forEachVector(t, "aggregate", func(t *testing.T, v vector) {
		var sigs []*Signature
		for _, b := range unhexList(t, v.Input) {
			sig, err := SignatureFromBytes(b)
			if err != nil {
				t.Fatal(err)
			}
			sigs = append(sigs, sig)
		}
		agg, err := Aggregate(sigs)
		if v.Output == nil {
			if err == nil {
				t.Fatal("expected error")
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		got := agg.Bytes()
		if want := unhex(t, v.Output); string(got[:]) != string(want) {
			t.Fatalf("aggregate mismatch:\n got  %x\n want %x", got, want)
		}
	})
}

func TestFastAggregateVerify(t *testing.T) {
	forEachVector(t, "fast_aggregate_verify", func(t *testing.T, v vector) {
		want := v.Output.(bool)
		pks, ok := parsePubkeys(t, unhexList(t, in(v)["pubkeys"]))
		sig, err := SignatureFromBytes(unhex(t, in(v)["signature"]))
		if !ok || err != nil {
			if want {
				t.Fatal("unexpected deserialization failure")
			}
			return
		}
		got, err := FastAggregateVerify(pks, unhex(t, in(v)["message"]), sig)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestAggregateVerify(t *testing.T) {
	forEachVector(t, "aggregate_verify", func(t *testing.T, v vector) {
		want := v.Output.(bool)
		pks, ok := parsePubkeys(t, unhexList(t, in(v)["pubkeys"]))
		sig, err := SignatureFromBytes(unhex(t, in(v)["signature"]))
		if !ok || err != nil {
			if want {
				t.Fatal("unexpected deserialization failure")
			}
			return
		}
		got, err := AggregateVerify(pks, unhexList(t, in(v)["messages"]), sig)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestDeserializationG1(t *testing.T) {
	forEachVector(t, "deserialization_G1", func(t *testing.T, v vector) {
		_, err := PublicKeyFromBytes(unhex(t, in(v)["pubkey"]))
		if ok := err == nil; ok != v.Output.(bool) {
			t.Fatalf("got ok=%v (err=%v), want %v", ok, err, v.Output)
		}
	})
}

func TestDeserializationG2(t *testing.T) {
	forEachVector(t, "deserialization_G2", func(t *testing.T, v vector) {
		_, err := SignatureFromBytes(unhex(t, in(v)["signature"]))
		if ok := err == nil; ok != v.Output.(bool) {
			t.Fatalf("got ok=%v (err=%v), want %v", ok, err, v.Output)
		}
	})
}

func TestKeyValidate(t *testing.T) {
	inf := make([]byte, 48)
	inf[0] = 0xc0
	if valid, err := KeyValidate(inf); valid || !errors.Is(err, ErrInfinityPubkey) {
		t.Fatalf("identity pubkey: got (%v, %v), want (false, ErrInfinityPubkey)", valid, err)
	}
	if valid, err := KeyValidate(inf[:47]); valid || err == nil {
		t.Fatalf("short encoding: got (%v, %v), want (false, some error)", valid, err)
	}

	sk, err := SecretKeyFromBytes(append(make([]byte, 31), 1))
	if err != nil {
		t.Fatal(err)
	}
	pk := sk.PublicKey().Bytes()
	// A valid key must come back with a nil error, not just valid==true.
	if valid, err := KeyValidate(pk[:]); !valid || err != nil {
		t.Fatalf("valid pubkey: got (%v, %v), want (true, nil)", valid, err)
	}
}

func TestDecoderIsCompressedOnly(t *testing.T) {
	mk := func(n int, first byte) []byte {
		b := make([]byte, n)
		b[0] = first
		return b
	}

	g1 := []struct {
		name string
		in   []byte
		want error
	}{
		{"compressed infinity (canonical)", mk(48, 0xc0), nil},
		{"uncompressed infinity, well-formed but not compressed", mk(96, 0x40), ErrInvalidLength},
		{"all-zero uncompressed", mk(96, 0x00), ErrInvalidLength},
		{"trailing byte after a compressed point", append(mk(48, 0xc0), 0x30), ErrInvalidLength},
		{"one byte short", mk(47, 0xc0), ErrInvalidLength},
		{"empty", nil, ErrInvalidLength},
	}
	for _, tc := range g1 {
		t.Run("G1/"+tc.name, func(t *testing.T) {
			if _, err := PublicKeyFromBytes(tc.in); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}

	g2 := []struct {
		name string
		in   []byte
		want error
	}{
		{"compressed infinity (canonical)", mk(96, 0xc0), nil},
		{"uncompressed infinity, well-formed but not compressed", mk(192, 0x40), ErrInvalidLength},
		{"all-zero uncompressed", mk(192, 0x00), ErrInvalidLength},
		{"trailing byte after a compressed point", append(mk(96, 0xc0), 0x30), ErrInvalidLength},
		{"one byte short", mk(95, 0xc0), ErrInvalidLength},
		{"empty", nil, ErrInvalidLength},
	}
	for _, tc := range g2 {
		t.Run("G2/"+tc.name, func(t *testing.T) {
			if _, err := SignatureFromBytes(tc.in); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestEthFastAggregateVerifyEmpty(t *testing.T) {
	inf := make([]byte, 96)
	inf[0] = 0xc0
	sig, err := SignatureFromBytes(inf)
	if err != nil {
		t.Fatal(err)
	}
	if ok, _ := EthFastAggregateVerify(nil, []byte("m"), sig); !ok {
		t.Fatal("empty pubkeys + infinity signature must be valid")
	}
	if ok, _ := FastAggregateVerify(nil, []byte("m"), sig); ok {
		t.Fatal("plain FastAggregateVerify must reject empty pubkeys")
	}
}

func FuzzKeyValidate(f *testing.F) {
	seedFromField(f, "deserialization_G1", "pubkey")

	f.Fuzz(func(t *testing.T, data []byte) {
		KeyValidate(data) // must never panic, whatever garbage comes in
	})
}

func FuzzPublicKeyFromBytes(f *testing.F) {
	seedFromField(f, "deserialization_G1", "pubkey")

	f.Fuzz(func(t *testing.T, data []byte) {
		pk, err := PublicKeyFromBytes(data)
		if err != nil {
			return
		}
		enc := pk.Bytes()
		pk2, err := PublicKeyFromBytes(enc[:])
		if err != nil {
			t.Fatalf("re-parsing own compressed encoding failed: %v", err)
		}
		if !pk.p.Equal(&pk2.p) {
			t.Fatalf("round-trip mismatch: parsed %x, re-encoded %x", data, enc)
		}
	})
}

func FuzzSignatureFromBytes(f *testing.F) {
	seedFromField(f, "deserialization_G2", "signature")

	f.Fuzz(func(t *testing.T, data []byte) {
		sig, err := SignatureFromBytes(data)
		if err != nil {
			return
		}
		enc := sig.Bytes()
		sig2, err := SignatureFromBytes(enc[:])
		if err != nil {
			t.Fatalf("re-parsing own compressed encoding failed: %v", err)
		}
		if !sig.p.Equal(&sig2.p) {
			t.Fatalf("round-trip mismatch: parsed %x, re-encoded %x", data, enc)
		}
	})
}

func FuzzVerify(f *testing.F) {
	vs, err := loadVectors(vectorsDir, "verify")
	if err != nil {
		f.Fatalf("load vectors: %v", err)
	}
	for _, nv := range vs {
		m := in(nv.v)
		pkHex, ok1 := m["pubkey"].(string)
		sigHex, ok2 := m["signature"].(string)
		msgHex, ok3 := m["message"].(string)
		if !ok1 || !ok2 || !ok3 {
			continue
		}
		pkB, err1 := hexBytes(pkHex)
		sigB, err2 := hexBytes(sigHex)
		msgB, err3 := hexBytes(msgHex)
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}
		f.Add(pkB, sigB, msgB)
	}

	f.Fuzz(func(t *testing.T, pkBytes, sigBytes, msg []byte) {
		pk, err := PublicKeyFromBytes(pkBytes)
		if err != nil {
			return
		}
		sig, err := SignatureFromBytes(sigBytes)
		if err != nil {
			return
		}
		_, _ = Verify(pk, msg, sig) // must not panic
	})
}

func FuzzFastAggregateVerify(f *testing.F) {
	vs, err := loadVectors(vectorsDir, "fast_aggregate_verify")
	if err != nil {
		f.Fatalf("load vectors: %v", err)
	}
	for _, nv := range vs {
		m := in(nv.v)
		pkList, ok1 := m["pubkeys"].([]any)
		sigHex, ok2 := m["signature"].(string)
		msgHex, ok3 := m["message"].(string)
		if !ok1 || !ok2 || !ok3 {
			continue
		}
		var blob []byte
		ok := true
		for _, e := range pkList {
			s, o := e.(string)
			if !o {
				ok = false
				break
			}
			b, err := hexBytes(s)
			if err != nil {
				ok = false
				break
			}
			blob = append(blob, b...)
		}
		sigB, err2 := hexBytes(sigHex)
		msgB, err3 := hexBytes(msgHex)
		if !ok || err2 != nil || err3 != nil {
			continue
		}
		f.Add(blob, sigB, msgB)
	}

	const pkSize = bls12381.SizeOfG1AffineCompressed
	f.Fuzz(func(t *testing.T, pkBlob, sigBytes, msg []byte) {
		sig, err := SignatureFromBytes(sigBytes)
		if err != nil {
			return
		}
		var pks []*PublicKey
		for len(pkBlob) >= pkSize {
			chunk := pkBlob[:pkSize]
			pkBlob = pkBlob[pkSize:]
			if pk, err := PublicKeyFromBytes(chunk); err == nil {
				pks = append(pks, pk)
			}
		}
		_, _ = FastAggregateVerify(pks, msg, sig)    // must not panic
		_, _ = EthFastAggregateVerify(pks, msg, sig) // must not panic
	})
}
