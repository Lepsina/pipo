// Package spectest reads ethereum/consensus-spec-tests and
// ethereum/bls12-381-tests vectors from one repository-level directory.
//
// The vectors are hundreds of megabytes, so they are not committed: the
// Makefile downloads a pinned version into spectests/, which is gitignored.
// Every package tests against that single copy instead of keeping its own,
// so there is one vector version for the whole project.
package spectest

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/snappy"
	"gopkg.in/yaml.v3"
)

// Versions of the vector releases this code is written against.
const (
	SpecVersion     = "v1.6.0-beta.0"
	BLSTestsVersion = "v0.1.2"
)

const DirEnv = "PIPO_SPECTESTS"

const dirName = "spectests"

func Dir(tb testing.TB) string {
	tb.Helper()

	dir := os.Getenv(DirEnv)
	if dir == "" {
		root, err := moduleRoot()
		if err != nil {
			tb.Fatalf("locating module root: %v", err)
		}
		dir = filepath.Join(root, dirName)
	}
	if _, err := os.Stat(dir); err != nil {
		tb.Skipf("spec vectors not found at %s; run `make spectests`", dir)
	}
	checkVersions(tb, dir)
	return dir
}

// Preset returns the directory of one fork's vectors, e.g.
// Preset(t, "mainnet", "fulu") for spectests/mainnet/fulu.
func Preset(tb testing.TB, preset, fork string) string {
	tb.Helper()
	return filepath.Join(Dir(tb), preset, fork)
}

// BLSTestsDir returns the root of the bls12-381-tests vectors.
func BLSTestsDir(tb testing.TB) string {
	tb.Helper()
	return filepath.Join(Dir(tb), "bls-12-381-tests")
}

// moduleRoot walks up from the working directory to the directory holding
// go.mod. Tests run in their own package directory, so the vector path cannot
// be a fixed relative path.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found above %s", dir)
		}
		dir = parent
	}
}

// checkVersions fails the test when spectests/VERSION disagrees with the
// versions above. Running against vectors from another release would make the
// results meaningless, so this is a failure and not a skip.
func checkVersions(tb testing.TB, dir string) {
	tb.Helper()

	path := filepath.Join(dir, "VERSION")
	raw, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("reading %s: %v (re-run `make spectests`)", path, err)
	}

	got := map[string]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		name, version, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok {
			continue
		}
		got[name] = strings.TrimSpace(version)
	}

	want := map[string]string{
		"consensus-spec-tests": SpecVersion,
		"bls12-381-tests":      BLSTestsVersion,
	}
	for name, wantVersion := range want {
		if got[name] != wantVersion {
			tb.Fatalf("%s in %s is %q, want %q; re-run `make spectests`",
				name, path, got[name], wantVersion)
		}
	}
}

// Runner walks the case directories under Dir and runs one subtest per case.
type Runner struct {
	// Dir is the directory holding the cases, at any nesting depth.
	Dir string
	// Skip maps a case name, as reported by Case.Name, to the reason that
	// case is skipped.
	Skip map[string]string
}

// Run walks the cases under dir with no skips.
func Run(t *testing.T, dir string, fn func(t *testing.T, c Case)) {
	t.Helper()
	Runner{Dir: dir}.Run(t, fn)
}

// Run runs fn once per case, in a subtest named after the case's path
// relative to r.Dir.
func (r Runner) Run(t *testing.T, fn func(t *testing.T, c Case)) {
	t.Helper()

	names, err := caseNames(r.Dir)
	if err != nil {
		t.Fatalf("walking %s: %v", r.Dir, err)
	}
	if len(names) == 0 {
		t.Fatalf("no cases found under %s", r.Dir)
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			if reason, ok := r.Skip[name]; ok {
				t.Skip(reason)
			}
			fn(t, Case{Name: name, dir: filepath.Join(r.Dir, filepath.FromSlash(name)), t: t})
		})
	}
}

// caseNames returns the slash-separated paths, relative to dir, of every
// directory below it that has no subdirectories of its own. The vector layout
// nests differently per test type (ssz_static/<Type>/ssz_random/case_0,
// single_merkle_proof/<Type>/<proof>), so the leaf directory is what
// identifies a case rather than a fixed depth.
func caseNames(dir string) ([]string, error) {
	var names []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				return nil
			}
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel != "." {
			names = append(names, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}

// Case is one spec-test case directory. Its accessors fail the test on a
// missing or malformed file: a vector that cannot be read is a broken
// download, not a property of the code under test.
type Case struct {
	// Name is the case path relative to the runner's directory.
	Name string

	dir string
	t   *testing.T
}

// Path returns the path of a file inside the case directory.
func (c Case) Path(name string) string {
	return filepath.Join(c.dir, name)
}

// SSZ returns the decompressed contents of <name>.ssz_snappy.
func (c Case) SSZ(name string) []byte {
	c.t.Helper()

	path := c.Path(name + ".ssz_snappy")
	compressed, err := os.ReadFile(path)
	if err != nil {
		c.t.Fatalf("reading %s: %v", path, err)
	}
	raw, err := snappy.Decode(nil, compressed)
	if err != nil {
		c.t.Fatalf("snappy decode %s: %v", path, err)
	}
	return raw
}

// YAML decodes <name>.yaml into out, which must be a pointer.
func (c Case) YAML(name string, out any) {
	c.t.Helper()

	path := c.Path(name + ".yaml")
	raw, err := os.ReadFile(path)
	if err != nil {
		c.t.Fatalf("reading %s: %v", path, err)
	}
	if err := yaml.Unmarshal(raw, out); err != nil {
		c.t.Fatalf("parsing %s: %v", path, err)
	}
}

// Root returns the expected hash tree root from roots.yaml, which ssz_static
// cases carry alongside the serialized object.
func (c Case) Root() [32]byte {
	c.t.Helper()

	var fixture struct {
		Root string `yaml:"root"`
	}
	c.YAML("roots", &fixture)
	return c.Bytes32(fixture.Root)
}

// Bytes32 parses a 0x-prefixed 32-byte hex string from a vector.
func (c Case) Bytes32(s string) [32]byte {
	c.t.Helper()

	raw, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		c.t.Fatalf("case %s: decoding %q: %v", c.Name, s, err)
	}
	if len(raw) != 32 {
		c.t.Fatalf("case %s: %q is %d bytes, want 32", c.Name, s, len(raw))
	}
	var out [32]byte
	copy(out[:], raw)
	return out
}
