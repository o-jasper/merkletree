package hash_extra

import "hash"

type HashResult []byte
type Hasher struct { hash.Hash }

func greater(x []byte, y []byte) bool {
	if len(x) != len(y) { 
		print(len(x), ",", len(y), "\n")
		panic("Unequal sized integer comparison (probably different hash interface types)")
	}
	for i := range x {
		if x[i] > y[i] { return true }
		if x[i] < y[i] { return false }
	}
	return true
}

func (hasher Hasher) rw(input []byte) {	// Reset, then write.
	hasher.Reset()
	hasher.Write(input)
}

func (hasher Hasher) S() HashResult { // Current sum.
	return hasher.Sum([]byte{})
}

func (hasher Hasher) H(input []byte) HashResult {
	hasher.rw(input)
	return hasher.S()
}
func (hasher Hasher) HwI(i uint64, data []byte) HashResult {
	hasher.Reset()
	for j := 0 ; j < 8 ; j ++ {
		hasher.Write([]byte{byte(i/72057594037927936)})
		i *= 256
	}
	hasher.Write(data)
	return hasher.S()
}

// Combine pair of hashes unorderedly.
func (hasher Hasher) H_U2(h1, h2 HashResult) HashResult {
	if greater(h1, h2) {
		hasher.rw(h1)
		hasher.Write(h2)
		return hasher.S()
	} else {
		hasher.rw(h2)
		hasher.Write(h1)
		return hasher.S()
	}
}
// .. orderedly.
func (hasher Hasher) H_2(a, b HashResult) HashResult {
	hasher.rw(a)
	hasher.Write(b)
	return hasher.S()
}

func ByteSliceEqual(a []byte, b []byte) bool {
	if len(a) != len(b) { return false }
	for i := range a { if a[i] != b[i] { return false } }
	return true
}

func HashEqual(a HashResult, b HashResult) bool {
	return ByteSliceEqual(a, b)
}

//Calculate expected root, given the path.
func (hasher Hasher) MerkleExpectedRoot(H_leaf HashResult, path []HashResult) HashResult {
	x := H_leaf
	for _, el := range path {	x = hasher.H_U2(el, x) }
	return x
}

//Checks a root.
func (hasher Hasher) MerkleVerifyH(root, Hleaf HashResult, path []HashResult) bool {
	return HashEqual(hasher.MerkleExpectedRoot(Hleaf, path), root)
}
func (hasher Hasher) MerkleVerify(root HashResult, leaf []byte, path []HashResult) bool {
	return hasher.MerkleVerifyH(root, hasher.H(leaf), path)
}

func (hasher Hasher) MerkleVerifyWithIndex(root HashResult, i uint64, leaf []byte, path []HashResult) bool {
	return hasher.MerkleVerifyH(root, hasher.HwI(i, leaf), path)
}
