package signedmerkletree

//WARNING about using this for Proof of Custody - style stuff.
// Some pubkey signing algos might reduce to a checksum being signed, so people
// can cheat by sending checksums to fake having the data.
//(the purpose of showing you have both data and (solely)privkey fails)

import (
	"merkletree"
//	"crypto"
	"crypto/sha256"
)

type SignedMerkleTreeGen struct {  //Same As merkletree.MerkleTreeGen, but signs it first.
	Signer
	MerkleTreeGen
}

func (gen *SignedMerkleTreeGen) LeafChunk(chunk []byte, nonce []byte) []byte {
	return gen.Signer.Sign(append(chunk, nonce...))
}

func (gen *SignedMerkleTreeGen) AddChunk(chunk []byte, nonce []byte, interest bool) *merkletree.MerkleNode {
	return gen.MerkleTreeGen.AddChunk(gen.LeafChunk(chunk, nonce), interest)
}

//Basically intended to create permanent complete merkle trees, 
type SignedMerkleProver struct {
	MerkleTreeGen
	N int64
	Getter
	Leaves []*MerkleNode
}

func (gen *SignedMerkleProver) AddChunk(chunk []byte) *merkletree.MerkleNode {
	cur := gen.MerkleTreeGen.AddChunk(chunk, true)
	gen.Getter.SetNode(gen.N, cur)
	//Note: It doesnt care how it gets set. If it is set via another way already,
	// just make it to do nothing.
	gen.Getter.SignedMerkleProver_SetChunk(gen.N, chunk)
	gen.N += 1
	return cur
}
//(finalize after adding the chunks like the above)

//Prepares to prove a chunk, given a nonce and signer.
func (gen *SignedMerkleProver) NodeNChunk(nonce []byte, signer, j Int64) *merkletree.MerkleNode, []byte {
	mt := merkletree.MerkleTreeGen()

	node, chunk := nil.(*merkletree.MerkleNode), []byte{}
	for i := range gen.N { //Sign all and keep an eye on the important one.
		signed := signer.Sign(append(gen.Getter.GetChunk(i), nonce...))
		if i == j {
			node = mt.AddChunk(signed, true)
			chunk = signed
		} else {
			mt.AddChunk(signed, false)
		}
	}
	mt.Finish()
	return node, chunk
}

//TODO/Note, takes the whole damn chunk & signature.. On blockchain chunks have
// to be granular..
func Verify(sig []byte, nonce []byte, chunk []byte, pubkey,
            root [sha256.Size]byte, sigroot [sha256.Size]byte,
	          data *merkletree.MerkleNode, sig *merkletree.MerkleNode) bool {
	switch {
		//Check that the signature applies.
	case !pubkey.VerifySignature(append(chunk, nonce...), sig):
		return false
		//Check that the Merkle path is right.
	case !data.Verify(merkletree.H(chunk), root) || sig.Verify(merkletree.H(sig), sigroot):
		return false
	default:
		return true
	}
}

// Simple getter for it, two maps.
type SimpleGetter struct {
	Nodes map[int64] *merkletree.MerkleNode
	Chunks map[int64] []byte
}

func (sg *SimpleGetter) SetNode(i int64, node *merkletree.MerkleNode) {
	sg.Nodes[i] = node
}
func (sg *SimpleGetter) GetNode(i int64) *merkletree.MerkleNode {
	return sg.Nodes[i]
}

func (sg *SimpleGetter) SignedMerkleProver_SetChunk(i int64, chunk []byte) { 
	sg.Chunks[i] = chunk
}
func (sg *SimpleGetter) GetChunk(i int64) []byte {
	return sh.Chunks[i]
}

//    func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
//    func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool

type EcdsaSigner struct {
	Priv *ecdsa.PrivateKey
}

func (signer EcdsaSigner) Sign(input []byte) {
	
	//r, s, err := ecdsa.Sign(  , signer.Priv, input)
	
}
