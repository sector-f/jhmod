package nvc

import (
	"fmt"
	"hash/fnv"
)

// TODO consider using the Stringer interface.
func Hash2String(hash uint64) string {
	return fmt.Sprintf("%016x", hash)
}

func String2Hash(s string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return hash.Sum64()
}
