package ssz

import (
	"crypto/sha256"
)

func hash(a, b []byte) [32]byte {
	buf := make([]byte, 0, len(a)+len(b))
	buf = append(buf, a...)
	buf = append(buf, b...)
	return sha256.Sum256(buf)
}
