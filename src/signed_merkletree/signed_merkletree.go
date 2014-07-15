package signed_merkletree

//WARNING about using this for Proof of Custody - style stuff.
// Some pubkey signing algos might reduce to a checksum being signed, so people
// can cheat by sending checksums to fake having the data.
//(the purpose of showing you have both data and (solely)privkey fails)
//
// So either the signature does not do that, or the chunks should be (nearly) as
// granular as to have ~ the size of a checksum.

import (
	"merkletree"
	"crypto/sha256"
)

type Getter interface {
	GetNode(int64) *merkletree.MerkleNode
	SetNode(int64, *merkletree.MerkleNode)

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
	merkletree.MerkleTreeGen
	N int64
	Getter
}

// Adds non-signed chunks.
func (gen *SignedMerkleProver) AddChunk(chunk []byte) *merkletree.MerkleNode {
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
func (gen *SignedMerkleProver) AddAllSigned(nonce []byte, signer Signer) (*merkletree.MerkleNode, *SignedMerkleProver) {
	smp := NewSignedMerkleProver()
	for i := int64(0) ; i < gen.N ; i++ {
		smp.AddChunk(signer.Sign(append(gen.Getter.GetChunk(i), nonce...)))
	}
	return smp.Finish(), &smp
}

//TODO/NOTE, takes the whole damn chunk & signature.. Or blockchain chunks have
// to be granular..
func Verify(sig []byte, chunk []byte, nonce []byte, pubkey Pubkey,
            root [sha256.Size]byte, sigroot [sha256.Size]byte,
            node *merkletree.MerkleNode, sig_node *merkletree.MerkleNode) bool {
	switch {
		//Check that the signature applies.
	case !pubkey.VerifySignature(sig, append(chunk, nonce...)):
		return false
		//Check that the Merkle path is right.
	case !node.Verify(merkletree.H(chunk), root) || sig_node.Verify(merkletree.H(sig), sigroot):
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
	return sg.Chunks[i]
}

func NewSimpleGetter() SimpleGetter {
	return SimpleGetter{map[int64]*merkletree.MerkleNode{}, map[int64][]byte{}}
}

func NewSignedMerkleProver() SignedMerkleProver {
	getter := NewSimpleGetter()
	return SignedMerkleProver{merkletree.NewMerkleTreeGen(), int64(0), &getter}
}
