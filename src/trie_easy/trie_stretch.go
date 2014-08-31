package trie_easy

import "fmt"

// --- Stretch with just one branch.

type TrieStretch struct {
	Stretch  []byte
	End      TrieNode
}

func (n *TrieStretch) Downward(str []byte, i int64) (*TrieNode, int64) {
	if i < 2*int64(len(str) - len(n.Stretch)) { // Range inside the stretch.
		return nil, i
	}
	for j := int64(0) ; j < int64(len(n.Stretch)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			return nil, i  // Back to begining.
		}
	}
	// Continue at end.
	return n.End.Downward(str, i + 2*int64(len(n.Stretch)))
}

func (n *TrieStretch) Get(str []byte, i int64) interface{} {
	if i < 2*int64(len(str) - len(n.Stretch)) { // Range inside the stretch.(nothing in there)
		return nil
	}
	return n.End.Get(str, i + 2*int64(len(n.Stretch)))
}

func (n* TrieStretch) SetRaw(str []byte, i int64, to interface{}) TrieNodeInterface {
	if i%2 != 0 { panic("TrieStretches should start at whole bytes.") }
	// The hard part.
	for j := int64(0) ; i/2 + j < int64(len(str)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			a1, a2 := str[i/2 + j]%16, str[i/2 + j]/16  //Added nibbles.
			g1, g2 := n.Stretch[j]%16, n.Stretch[j]/16  //Already existing nibble route.

			fmt.Println("F", n.Stretch, i, j, str)

			// Two nodes(stretches start even)
			first, second := NewTrieNode16(nil), NewTrieNode16(nil)
			// Connect them.
			first.Sub[g1] = NewTrieNode(second)
			// Connect to what is after.
			if j < 2*int64(len(str)) { // It is a bit of stretch.
				fmt.Println("-", n.Stretch[j+1:], g1+16*g2)
				after_stretch := &TrieStretch{ Stretch : n.Stretch[j+1:], End : n.End }
				second.Sub[g2] = NewTrieNode(after_stretch)
			} else{  // It is the current end.
				second.Sub[g2] = n.End
			}

			if a1 != g1 { // Breaks out of first one.
				first.Sub[a1].SetI(str, i + 2*j, to)
			} else if a2 != g2 { // Breaks out of first one.
				if a1 != g1 { panic("BUG") }  // (could be both, first goes)
				second.Sub[a2].SetI(str, i + 2*j + 1, to)
			} else { panic("BUG") }
			
			if j == 0 {
				return first
			} else {  // Prepend what was before.
				fmt.Println("-", n.Stretch[:j])
				return &TrieStretch{Stretch : n.Stretch[:j], End : NewTrieNode(first)}
			}
			panic("Unreachable")
		}
	}
	panic("BUG Didnt go downward properly.(if it split it should also have cut out now.)")
	return nil
}

type CreatorStretch struct{}

func (_ CreatorStretch) Extend(str []byte, i int64, final TrieNodeInterface) TrieNodeInterface {
	if i/2 + 1 == int64(len(str)) {
		if i%2 == 0 {
			first := NewTrieNode16(nil)
			first.Sub[str[i/2]%16].Actual = final
			return first
		} else {
			return final
		}
	}
	first := &TrieStretch{Stretch : str[i/2 + 1:], End : NewTrieNode(final)}
	if i%2 == 0 { // This one is even, so the next one isnt!
		m := NewTrieNode16(nil)
		m.Sub[str[i/2]%16].Actual = first
		return m  // This one as first, instead of other one.
	}
	return first
}
