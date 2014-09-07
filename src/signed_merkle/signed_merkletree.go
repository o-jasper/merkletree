package signed_merkle

//WARNING about using this for Proof of Custody - style stuff.
// Some pubkey signing algos might reduce to a checksum being signed, so people
// can cheat by sending checksums to fake having the data.
//(the purpose of showing you have both data and (solely)privkey fails)
//
// So either the signature does not do that, or the chunks should be (nearly) as
// granular as to have ~ the size of a checksum.

import (
	"merkle"
	"hash"
)

type Getter interface {
	GetNode(int64) *merkle.MerkleNode
	SetNode(int64, *merkle.MerkleNode)

	GetChunk(int64) []byte
	SignedMerkleProver_SetChunk(int64, []byte)
}

type Signer interface {
	Sign(input []byte) []byte
}
type Pubkey interface {
	VerifySignature(sig []byte, data []byte) bool
}

// Basically intended to create permanent complete merkle trees, 
type SignedMerkleProver struct {
	merkle.MerkleTreeGen
	N int64
	Getter
}

// Adds non-signed chunks.
func (gen *SignedMerkleProver) AddChunk(chunk []byte) *merkle.MerkleNode {
	cur := gen.MerkleTreeGen.AddChunk(chunk, true)
	gen.Getter.SetNode(gen.N, cur)
	// Note: It doesnt care how it gets set. If it is set via another way already,
	//  just make it to do nothing.
	gen.Getter.SignedMerkleProver_SetChunk(gen.N, chunk)
	gen.N += 1
	return cur
}
// (finalize after adding the chunks like the above)

// Prepares to prove a chunk, given a nonce and signer.
func (gen *SignedMerkleProver) AddAllSigned(nonce []byte, signer Signer) (*merkle.MerkleNode, *SignedMerkleProver) {
	smp := NewSignedMerkleProver(gen.Hash, gen.IncludeIndex)
	for smp.N < gen.N {
		smp.AddChunk(signer.Sign(append(gen.Getter.GetChunk(smp.N), nonce...)))
	}
	return smp.Finish(), &smp
}

// Note: assumes the index and nonce is already at the verifier.
type SignedMerkleProof struct {
	node      *merkle.MerkleNode
	sig_node  *merkle.MerkleNode
	chunk     []byte
	sig_chunk []byte
}

func (gen *SignedMerkleProver) NewSignedMerkleProof_FromIndex(signed *SignedMerkleProver, index int64) SignedMerkleProof {
	return SignedMerkleProof{ 
		node      : gen.Getter.GetNode(index),
		sig_node  : signed.Getter.GetNode(index),
		chunk     : gen.Getter.GetChunk(index),
		sig_chunk : signed.Getter.GetChunk(index) }
}

//TODO/NOTE, takes the whole damn chunk & signature.. Or blockchain chunks have
// to be granular..
func (proof *SignedMerkleProof) Verify(nonce []byte, pubkey Pubkey, root hash.Hash, sig_root hash.Hash) int8 {
	//Check that the signature applies.
	if !pubkey.VerifySignature(proof.sig_chunk, append(proof.chunk, nonce...)) {
		return merkle.WrongSig
	} else { //Check that the Merkle paths are right.
		if r := proof.sig_node.Verify(sig_root, proof.sig_chunk) ; r == merkle.Correct {
			return proof.node.Verify(root, proof.chunk) //It takes over.
		} else { //Makes it the Signature path error version.
			return r + merkle.Merkletree_NPathWrongs
		}
	}
}

//Turns the proof into bytes.
/*func (proof *SignedMerkleProof) ByteProof() {
	ret := append(proof.sig_chunk, proof.chunk...)
	ret = append(ret, byte(len(proof.node.ByteProof())/32))
	ret = append(ret, proof.node.ByteProof()...) //Issue.. it is [][32]byte
	ret = append(ret, proof.sig_node.ByteProof()...)
	return ret
}*/

// Simple getter for it, two maps.
type SimpleGetter struct {
	Nodes map[int64] *merkle.MerkleNode
	Chunks map[int64] []byte
}

func (sg *SimpleGetter) SetNode(i int64, node *merkle.MerkleNode) {
	sg.Nodes[i] = node
}
func (sg *SimpleGetter) GetNode(i int64) *merkle.MerkleNode {
	return sg.Nodes[i]
}

func (sg *SimpleGetter) SignedMerkleProver_SetChunk(i int64, chunk []byte) { 
	sg.Chunks[i] = chunk
}
func (sg *SimpleGetter) GetChunk(i int64) []byte {
	return sg.Chunks[i]
}

func NewSimpleGetter() SimpleGetter {
	return SimpleGetter{map[int64]*merkle.MerkleNode{}, map[int64][]byte{}}
}

func NewSignedMerkleProver(h hash.Hash, include_index bool) SignedMerkleProver {
	getter := NewSimpleGetter()
	return SignedMerkleProver{merkle.NewMerkleTreeGen(h, include_index), int64(0), &getter}
}
