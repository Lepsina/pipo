package bls

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const vectorsDir = "../testdata/bls-12-381-tests"

type vector struct {
	Input  any `yaml:"input"`
	Output any `yaml:"output"`
}

func forEachVector(t *testing.T, group string, fn func(t *testing.T, v vector)) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(vectorsDir, group, "*.yaml"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no vectors found for %s: %v", group, err)
	}
	for _, f := range files {
		t.Run(strings.TrimSuffix(filepath.Base(f), ".yaml"), func(t *testing.T) {
			raw, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			var v vector
			if err := yaml.Unmarshal(raw, &v); err != nil {
				t.Fatal(err)
			}
			fn(t, v)
		})
	}
}

func in(v vector) map[string]any { return v.Input.(map[string]any) }

func unhex(t *testing.T, x any) []byte {
	t.Helper()
	s, ok := x.(string)
	if !ok {
		t.Fatalf("expected hex string, got %T", x)
	}
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		t.Fatal(err)
	}
	return b
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

// pubkeys/signature that fail to deserialize make the whole check invalid.
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
	if KeyValidate(inf) {
		t.Fatal("identity pubkey must be invalid")
	}
	if KeyValidate(inf[:47]) {
		t.Fatal("short encoding must be invalid")
	}
	sk, _ := SecretKeyFromBytes(append(make([]byte, 31), 1))
	pk := sk.PublicKey().Bytes()
	if !KeyValidate(pk[:]) {
		t.Fatal("valid pubkey rejected")
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
