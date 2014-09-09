package merkle_common

import (
	"encoding/gob"
	"bytes"
)

func GetBytes(key interface{}) []byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(key)
	return buf.Bytes()
}
