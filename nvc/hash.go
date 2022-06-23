package nvc

import (
	"fmt"
	"hash/fnv"
)

func Hash2String(hash Hash) string {
	return fmt.Sprintf("%016x", uint64(hash))
}

func String2Hash(s string) Hash {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return Hash(hash.Sum64())
}

func (h Hash) String() string {
	return Hash2String(h)
}
