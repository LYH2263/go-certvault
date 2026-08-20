package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
)

var seq uint64

// New 生成可读唯一 ID。
func New(prefix string) string {
	n := atomic.AddUint64(&seq, 1)
	var b [4]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s-%d-%s", prefix, n, hex.EncodeToString(b[:]))
}
