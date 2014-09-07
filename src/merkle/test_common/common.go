package test_common

import (
	"math/rand"
	"encoding/hex"
	"hash"
)

//Generating chunks.
func Rand_bytes(r *rand.Rand, n int64) []byte {
	out := []byte{}
	for i := int64(0) ; i < n ; i++ {
		out = append(out, byte(r.Int63n(256)))
	}
	return out
}

func Rand_range(r *rand.Rand, fr int64, to int64) int64 {
	return fr + r.Int63n(to - fr)
}

func Rand_chunk(r *rand.Rand, n_min int64, n_max int64) []byte {
	return Rand_bytes(r, Rand_range(r, n_min, n_max))
}

func HashBytes(h hash.Hash) []byte { return h.Sum([]byte{}) }
func HashStr(h hash.Hash) string { return hex.EncodeToString(HashBytes(h)) }
