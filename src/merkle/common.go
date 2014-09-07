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

func greater(x []byte, y []byte) bool {
	if len(x) != len(y) { 
		panic("Unequal sized integer comparison (probably different hash interface types)")
	}
	for i := range x {
		if x[i] > y[i] { return true }
		if x[i] < y[i] { return false }
	}
	return true
}

func H_2(h1 hash.Hash, h2 hash.Hash) hash.Hash {
	d1, d2 := h1.Sum([]byte{}), h2.Sum([]byte{})
	h_out := h1  // TODO ... wtf how do you make a 'same hash.hash'...
	h_out.Reset()
	if greater(d1, d2) {
		h_out.Write(d1)
		h_out.Write(d2)
	} else {
		h_out.Write(d2)
		h_out.Write(d1)
	}
	return h_out
}

func byteSliceEqual(a []byte, b []byte) bool {
	if len(a) != len(b) { return false }
	for i := range a { if a[i] != b[i] { return false } }
	return true
}

func HashEqual(a hash.Hash, b hash.Hash) bool {
	return byteSliceEqual(a.Sum([]byte{}), b.Sum([]byte{}))
}
