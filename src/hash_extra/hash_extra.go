package hash_extra

import (
	"hash"
)

// Allows you to create a hash.Hash with just a result in it.
//
// The hash standard library fails to deliver a bit on this area..

type HashResult struct {
	Result []byte
	Continue hash.Hash
}

//(for easier path converting)
func HashResultList(input [][]byte, cont hash.Hash) []hash.Hash {
	ret := []hash.Hash{}
	for _, el := range input { ret = append(ret, HashResult{el, cont}) }
	return ret
}

func (hr HashResult) Sum(b []byte) []byte {
	if len(b) != 0 { panic("This is a dud hash.Hash intended only for result!") }
	return hr.Result
}

func (hr HashResult) Reset() { panic("This is a dud hash.Hash") }
func (hr HashResult) Size() int { return len(hr.Result) }
func (hr HashResult) BlockSize() int { panic("This is a dud hash.Hash") }
func (hr HashResult) Write(_ []byte) (_ int, _ error){ panic("This is a dud hash.Hash") }

func ContinueUse(input hash.Hash) hash.Hash {
	if hr, yes := input.(HashResult) ; yes { return hr.Continue }
	if input == nil { panic("Problem, dont have a hash!") }
	input.Reset()
	return input
}

//
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

// Combine pair of hashes unorderedly.
func H_U2(h1, h2 hash.Hash) hash.Hash {
	d1, d2 := h1.Sum([]byte{}), h2.Sum([]byte{})
	h_out := ContinueUse(h1)
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

func H_2(a, b hash.Hash) hash.Hash {
	h := ContinueUse(a)
	h.Write(a.Sum([]byte{}))
	h.Write(b.Sum([]byte{}))
	return h
}
