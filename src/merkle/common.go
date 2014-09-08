package merkle

import (
	"hash"
	"encoding/gob"
	"bytes"
)

func getBytes(key interface{}) []byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(key)
	return buf.Bytes()
}

func byteSliceEqual(a []byte, b []byte) bool {
	if len(a) != len(b) { return false }
	for i := range a { if a[i] != b[i] { return false } }
	return true
}

func HashEqual(a hash.Hash, b hash.Hash) bool {
	return byteSliceEqual(a.Sum([]byte{}), b.Sum([]byte{}))
}
