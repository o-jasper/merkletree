package data_n_privkey

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
type PoCProver struct {
	MerkleTreeGen
	N int64
	Getter
	Leaves []*MerkleNode
}

func (gen *PoCProver) AddChunk(chunk []byte) *merkletree.MerkleNode {
	cur := gen.MerkleTreeGen.AddChunk(chunk, true)
	gen.Getter.SetNode(gen.N, cur)
	//Note: It doesnt care how it gets set. If it is set via another way already,
	// just make it to do nothing.
	gen.Getter.PoCProveSetChunk(gen.N, chunk)
	gen.N += 1
	return cur
}


//(finalize after adding the chunks like the above)

//Prepares to prove a chunk, given a nonce and signer.
func (gen *PoCProver) NodeNChunk(nonce []byte, signer, j Int64) *merkletree.MerkleNode, []byte {
	mt := merkletree.MerkleTreeGen()

	node, chunk := nil.(*merkletree.MerkleNode), []byte{}
	for i := range gen.N {
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

//Items needed to prove:
// * Regular Merkle path.
// * Nonced&signed-leafs - Merkle path.
// * Signed leaf
//
//Actions:
// * Check the paths.
//   + The steps
//   + Correspondence to roots.
//   + Correspondence to leaves.
// * Check the signature.


/*
func (interest *DataNPrivKeyInterest) IsValid(recurse int32) *DataNPrivKeyInterest, bool {
	s, svalid := interest.Sigs.IsValid(recurse)
	d, dvalid := interest.Data.IsValid(recurse)  // NOTE crosscheck left/right paths sequence?
	return &DataNPrivKeyInterest{Sigs:s, Data:d}, svalid && dvalid
}

func (interest *DataNPrivKeyInterest) CorrespondsToPubkey(pubkey interface{}) bool {
	return pubkey.Verify(interest.ChunkSig, interest.Chunk)
}

func (interest *DataNPrivKeyInterest) Verify(Hts [sha256.Size]byte, Htd [sha256.Size]byte, pubkey interface) bool {
	root, internal := interest.IsValid(-1)
	return (internal && 
		interest.CorrespondsToPubkey(pubkey) && 
		root.Sigs.CorrespondsToHash(Hts) && 
		root.Data.CorrespondsToHash(Htd) &&
		interest.Sig.CorrespondsToChunk(Interest.ChunkSig) &&
		interest.Data.CorrespondsToChunk(Interest.Chunk))
}
*/

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

func (sg *SimpleGetter) PoCProveSetChunk(i int64, chunk []byte) { 
	sg.Chunks[i] = chunk
}
func (sg *SimpleGetter) GetChunk(i int64) []byte {
	return sh.Chunks[i]
}
