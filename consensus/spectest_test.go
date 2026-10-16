package consensus

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/Lepsina/pipo/internal/spectest"
	"github.com/Lepsina/pipo/internal/ssz"
)

// sszType is what a spec-test fixture needs from a type under test:
// decode from wire bytes, re-encode back to wire bytes, and compute the
// root to compare against the fixture's expected root.
type sszType interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashTreeRooter
}

// specTestTypes maps a consensus-spec-tests ssz_static type name to a
// constructor for the matching Go type.
var specTestTypes = map[string]func() sszType{
	"BeaconBlockHeader":           func() sszType { return new(BeaconBlockHeader) },
	"SyncAggregate":               func() sszType { return new(SyncAggregate) },
	"SyncCommittee":               func() sszType { return new(SyncCommittee) },
	"ExecutionPayloadHeader":      func() sszType { return new(ExecutionPayloadHeader) },
	"LightClientHeader":           func() sszType { return new(LightClientHeader) },
	"LightClientBootstrap":        func() sszType { return new(LightClientBootstrap) },
	"LightClientUpdate":           func() sszType { return new(LightClientUpdate) },
	"LightClientFinalityUpdate":   func() sszType { return new(LightClientFinalityUpdate) },
	"LightClientOptimisticUpdate": func() sszType { return new(LightClientOptimisticUpdate) },
}

// TestSpecSSZStatic checks every type in specTestTypes against the mainnet
// Fulu ssz_static vectors: decode, root, and byte-identical re-encode.
func TestSpecSSZStatic(t *testing.T) {
	forkDir := spectest.Preset(t, "mainnet", "fulu")

	for typeName, newInstance := range specTestTypes {
		t.Run(typeName, func(t *testing.T) {
			dir := filepath.Join(forkDir, "ssz_static", typeName)
			spectest.Run(t, dir, func(t *testing.T, c spectest.Case) {
				wantBytes := c.SSZ("serialized")
				wantRoot := c.Root()

				v := newInstance()
				if err := v.Unmarshal(wantBytes); err != nil {
					t.Fatalf("Unmarshal: %v", err)
				}

				if gotRoot := v.HashTreeRoot(); gotRoot != wantRoot {
					t.Fatalf("HashTreeRoot mismatch: got %x, want %x", gotRoot, wantRoot)
				}

				gotBytes, err := v.Marshal()
				if err != nil {
					t.Fatalf("Marshal: %v", err)
				}
				if !bytes.Equal(gotBytes, wantBytes) {
					t.Fatalf("Marshal round trip mismatch:\ngot  %x\nwant %x", gotBytes, wantBytes)
				}
			})
		})
	}
}
