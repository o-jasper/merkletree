package merkle

import (
	"encoding/gob"
	"bytes"
)

func getBytes(key interface{}) []byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(key)
	return buf.Bytes()
}
