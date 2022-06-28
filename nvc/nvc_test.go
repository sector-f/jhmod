package nvc

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"

	"github.com/dsnet/golib/memfile"
)

type file struct {
	name     string
	contents []byte
}

func makeTestNVC(t *testing.T, compress bool, files ...file) *memfile.File {
	nvcFile := &memfile.File{}
	writer, err := NewWriter(nvcFile, uint32(len(files)))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		pr, pw := io.Pipe()
		go func() {
			pw.Write(f.contents)
			pw.Close()
		}()

		if compress {
			_, err = writer.CreateCompressed(pr, String2Hash(f.name), zlib.DefaultCompression)
		} else {
			_, err = writer.Create(pr, String2Hash(f.name))
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	err = writer.Finalize()
	if err != nil {
		t.Fatal(err)
	}

	nvcFile.Seek(0, io.SeekStart)
	return nvcFile
}

func TestReadFile(t *testing.T) {
	input := "Testing\n"

	nvcFile := makeTestNVC(t, false, file{"/path/to/file", []byte(input)})
	parsed, err := Parse(nvcFile)
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

func TestReadMultiFile(t *testing.T) {
	files := []file{
		{"foo", []byte("foobar\n")},
		{"fox.txt", []byte("The quick brown fox jumps over the lazy dog\n")},
	}

	nvcFile := makeTestNVC(t, false, files...)
	parsed, err := Parse(nvcFile)
	if err != nil {
		t.Fatal(err)
	}

	for i, f := range files {
		parsedContents, err := parsed.File(parsed.EntryOrder[i])
		if err != nil {
			t.Fatal(err)
		}

		expected := f.contents
		if bytes.Compare(parsedContents, expected) != 0 {
			t.Fatalf("Got %v, expected %v\n", parsedContents, expected)
		}
	}
}

func TestReadCompressed(t *testing.T) {
	files := []file{
		{"foo", []byte("foobar\n")},
		{"fox.txt", []byte("The quick brown fox jumps over the lazy dog\n")},
	}

	nvcFileUncompressed := makeTestNVC(t, false, files...)
	parsedUncompressed, err := Parse(nvcFileUncompressed)
	if err != nil {
		t.Fatal(err)
	}

	nvcFileCompressed := makeTestNVC(t, true, files...)
	parsedCompressed, err := Parse(nvcFileCompressed)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that different data was written by the Writer
	if bytes.Compare(nvcFileUncompressed.Bytes(), nvcFileCompressed.Bytes()) == 0 {
		t.Fatal("Compressed and uncompressed files are the same")
	}

	// Verify that same data is returned after parsing
	for i := range files {
		parsedContentsUncompressed, err := parsedUncompressed.File(parsedUncompressed.EntryOrder[i])
		if err != nil {
			t.Fatal(err)
		}

		parsedContentsCompressed, err := parsedCompressed.File(parsedCompressed.EntryOrder[i])
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Compare(parsedContentsUncompressed, parsedContentsCompressed) != 0 {
			t.Fatalf("%v does not match %v\n", parsedContentsUncompressed, parsedContentsCompressed)
		}
	}

}
