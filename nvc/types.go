package nvc

type EntryFlags uint32
type Hash uint64

const (
	Uncompressed EntryFlags = 0
	Compressed   EntryFlags = 1
	Encrypted    EntryFlags = 3

	magic = "nvc1d\x00\x00\x00"
)

type TocEntry struct {
	Hash        Hash
	Offset      uint32
	RawLength   uint32 // Uncompressed length
	Length      uint32 // Actual length in file
	Flags       EntryFlags
}
