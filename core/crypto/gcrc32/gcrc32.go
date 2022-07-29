package gcrc32

import (
	"github.com/Ravior/gserver/core/util/gconv"
	"hash/crc32"
)

// Encrypt encrypts any type of variable using CRC32 algorithms.
// It uses gconv package to convert <v> to its bytes type.
func Encrypt(v interface{}) uint32 {
	return crc32.ChecksumIEEE(gconv.Bytes(v))
}
