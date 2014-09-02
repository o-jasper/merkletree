package trie_easy

import "fmt"

// --- Stretch with just one branch.

type TrieStretch struct {
	Stretch  []byte
	End      Trie
}


func (n *TrieStretch) Down1(str []byte, i int64, _ bool) *Trie {
	if i + 2*int64(len(n.Stretch)) < 2*int64(len(str)) { // Range inside the stretch.
		return nil
	}
	for j := int64(0) ; j < int64(len(n.Stretch)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			return nil  // Back to begining.
		}
	}
	return &n.End
}

func (n *TrieStretch) Get(str []byte, i int64) interface{} {
	// Range inside the stretch.(nothing in there)
	if i + 2*int64(len(n.Stretch)) > 2*int64(len(str)) { return nil }
	if i < 2*int64(len(str)) { panic("String to short?!") }
	panic("Should be in next node")
	//return n.End.Get(str, i + 2*int64(len(n.Stretch)))
}

func (n* TrieStretch) SetRaw(str []byte, i int64, to interface{}) TrieInterface {
	if i%2 != 0 { panic("TrieStretches should start at whole bytes.") }
	// The hard part.
	for j := int64(0) ; j < int64(len(n.Stretch)) && i/2 + j < int64(len(str)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			a1, a2 := str[i/2 + j]%16, str[i/2 + j]/16  //Added nibbles.
			g1, g2 := n.Stretch[j]%16, n.Stretch[j]/16  //Already existing nibble route.

			fmt.Println("F", n.Stretch, i, j, str)

			// Two nodes(stretches start even)
			first, second := NewNode16(nil), NewNode16(nil)
			// Connect them.
			first.Sub[g1] = NewTrie(second)
			// Connect to what is after.
			if j < int64(len(n.Stretch) - 1) { // If is before the last, need stretch to end.
				fmt.Println("a-", n.Stretch[j+1:], g1+16*g2)
				after_stretch := &TrieStretch{ Stretch : n.Stretch[j+1:], End : n.End }
				second.Sub[g2] = NewTrie(after_stretch)
			} else{  // It is the last one, can go straight there.
				second.Sub[g2] = n.End
			}

			if a1 != g1 { // Breaks out of first one.
				fmt.Println("A", str[i/2 + j - 1:], i, j)
				first.Sub[a1].SetI(str, i + 2*j, to)
			} else if a2 != g2 { // Breaks out of first one.
				if a1 != g1 { panic("BUG") }  // (could be both, first goes)
				second.Sub[a2].SetI(str, i + 2*j + 1, to)
			} else { panic("BUG") }
			
			if j == 0 {
				return first
			} else {  // Prepend what was before.
				fmt.Println("b-", n.Stretch[:j])
				return &TrieStretch{Stretch : n.Stretch[:j], End : NewTrie(first)}
			}
			panic("Unreachable")
		}
	}
	panic("BUG Didnt go downward properly.(if it split it should also have cut out now.)")
	return nil
}

func (n* TrieStretch) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
//	if odd { panic("Stretches must be on even.") }
	return n.End.Actual.(TrieInterface).MapAll(data, append(pre, n.Stretch...), false, fun)
}

type CreatorStretch struct{}

func (_ CreatorStretch) Extend(str []byte, i int64, final TrieInterface) interface{} {
	if i/2 + 1 == int64(len(str)) {
		if i%2 != 0 { panic("if i == 2*len(str), shouldnt be here") }
		first := NewNode16(nil)
		first.Sub[str[i/2]%16].Actual = final
		return first
	}
	first := &TrieStretch{Stretch : str[i/2:], End : NewTrie(final)}
	if i%2 != 0 { // Not even, need one in the middle.
		m := NewNode16(nil)
		m.Sub[str[i/2]/16].Actual = first
		return m  // This one as first, instead of other one.
	}
	return first
}
