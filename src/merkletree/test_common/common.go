package test_common

import (
	"math/rand"
	"encoding/hex"
	"hash"
)

//Generating chunks.
func Rand_bytes(r *rand.Rand, n int32) []byte {
	out := []byte{}
	for i := int32(0) ; i < n ; i++ {
		out = append(out, byte(r.Int63n(256)))
	}
	return out
}

func Rand_range(r *rand.Rand, fr int32, to int32) int32 {
	return fr + r.Int31n(to - fr)
}

func Rand_chunk(r *rand.Rand, n_min int32, n_max int32) []byte {
	return Rand_bytes(r, Rand_range(r, n_min, n_max))
}

func HashBytes(h hash.Hash) []byte { return h.Sum([]byte{}) }
func HashStr(h hash.Hash) string { return hex.EncodeToString(HashBytes(h)) }
