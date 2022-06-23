package nvc

type EntryFlags uint32

const (
	UNCOMPRESSED EntryFlags = 0
	COMPRESSED   EntryFlags = 1
	ENCRYPTED    EntryFlags = 2

	MAGIC = "nvc1d\x00\x00\x00"
)

type TocEntry struct {
	Hash        uint64
	Offset      uint32
	RawLength   uint32 // Uncompressed length
	Length      uint32 // Actual length in file
	Flags       EntryFlags
}
