package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

/*
NewID ...
*/
func NewID() string {
	tb := make([]byte, 64)
	tc := binary.PutUvarint(tb, uint64(time.Now().UnixNano()))
	rb := make([]byte, 64)
	rc := binary.PutUvarint(rb, uint64(rand.Int63()))
	b := append(tb[:tc], rb[:rc]...)
	return fmt.Sprintf("%x", md5.Sum(b))
}
