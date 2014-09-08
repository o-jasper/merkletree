package trie_merkle

import (
	"encoding/gob"
	"bytes"
)

// https://stackoverflow.com/questions/23003793/convert-arbitrary-golang-interface-to-byte-array
func getBytes(key interface{}) []byte {
	if key == nil { return []byte{} }  // In this case it wusses out.
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(key)
	return buf.Bytes()
}
