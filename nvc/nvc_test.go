package nvc

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func makeTestNVC(file []byte) []byte {
	entry := TocEntry{
		Hash:      String2Hash("/path/to/file"),
		Offset:    8 + 4 + 24, // Magic bytes + ToC entry count + single ToC entry
		RawLength: uint32(len(file)),
		Length:    uint32(len(file)),
		Flags:     EntryFlagNoCompression,
	}

	nvcFile := bytes.NewBuffer([]byte(magic))
	binary.Write(nvcFile, binary.LittleEndian, uint32(1))
	binary.Write(nvcFile, binary.LittleEndian, entry)
	nvcFile.Write(file)

	return nvcFile.Bytes()
}

func TestReadUncompressedFile(t *testing.T) {
	input := "Testing\n"

	nvcFileBytes := bytes.NewReader(makeTestNVC([]byte(input)))
	parsed, err := Parse(nvcFileBytes)
	if err != nil {
		t.Fatal(err)
	}

	parsedContents, err := parsed.File(parsed.EntryOrder[0])
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(parsedContents, []byte(input)) != 0 {
		t.Fatalf("Got %v, expected %v\n", parsedContents, input)
	}
}
