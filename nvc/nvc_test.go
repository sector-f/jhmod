package nvc

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"

	"github.com/dsnet/golib/memfile"
)

func makeTestNVC(file []byte) []byte {
	nvcFile := &memfile.File{}
	writer, err := NewWriter(nvcFile, 1)
	if err != nil {
		panic(err)
	}

	pr, pw := io.Pipe()
	go func() {
		pw.Write(file)
		pw.Close()
	}()

	_, err = writer.CreateCompressed(pr, String2Hash("/path/to/file"), zlib.DefaultCompression)
	if err != nil {
		panic(err)
	}

	err = writer.Finalize()
	if err != nil {
		panic(err)
	}

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
