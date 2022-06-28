// Package nvc implements a parser for nvc archive files.
//
// The format of an archive is described below. All header values are stored on disk as little endian.
//
//   1. NVC magic bytes (8 bytes)
//   2. Number of entries in the archive's table of contents (32-bit little endian unsigned integer)
//   3. Table of Contents entries (24 * n bytes, where n is the previously-specified number of entries)
package nvc

import (
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
)

// EntryFlags describes how a file is stored on disk
type EntryFlags uint32

const (
	// EntryFlagNoCompression indicates that the file is stored in the NVC archive uncompressed
	EntryFlagNoCompression EntryFlags = 0
	// EntryFlagZlibCompression indicates that the file is stored in the NVC archive with zlib compression
	EntryFlagZlibCompression EntryFlags = 1
	// EntryFlagEncrypted indicates that the file is encrypted
	EntryFlagEncrypted EntryFlags = 3

	magic       = "nvc1d\x00\x00\x00" // NVC file type magic bytes
	headerLen   = 8 + 4               // Length of magic bytes + ToC entry count
	tocEntryLen = 24                  // Length of single ToC entry
)

var ErrNoMagicFound error = errors.New("nvc magic bytes not found")

type Archive struct {
	Entries    map[Hash]TocEntry // Map of hashes to table of contents entries
	EntryOrder []Hash            // List of entry hashes in the order that they are stored in the archive

	r io.ReadSeeker
}

// Parse reads r and attempts to interpret is as an NVC archive.
// This function takes ownership of r; it should not be used by the caller after Parse has been called.
// The returned Archive should not be used when the returned error is non-nil.
func Parse(r io.ReadSeeker) (Archive, error) {
	if magicErr := readMagic(r); magicErr != nil {
		return Archive{}, magicErr
	}

	count, countErr := readCount(r)
	if countErr != nil {
		return Archive{}, countErr
	}

	entries := make(map[Hash]TocEntry)
	order := make([]Hash, count)

	var i uint32
	for i = 0; i < count; i++ {
		entry, eErr := readEntry(r)
		if eErr != nil {
			return Archive{}, eErr
		}

		entries[entry.Hash] = entry
		order[i] = entry.Hash
	}

	a := Archive{
		Entries:    entries,
		EntryOrder: order,
		r:          r,
	}

	return a, nil
}

// File returns the data for the file that is reference by hash
func (a Archive) File(hash Hash) ([]byte, error) {
	entry, exists := a.Entries[hash]
	if !exists {
		return nil, errors.New("hash not present in archive")
	}

	_, err := a.r.Seek(int64(entry.Offset), 0)
	if err != nil {
		return nil, err
	}

	var reader io.Reader = a.r

	switch entry.Flags {
	case EntryFlagNoCompression:
		// Do nothing
	case EntryFlagZlibCompression:
		reader, err = zlib.NewReader(a.r)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported entry flag")
	}

	data := make([]byte, entry.RawLength)
	readBytes, err := io.ReadFull(reader, data)
	if err != nil {
		return nil, fmt.Errorf("error reading data after %d bytes: %w", readBytes, err)
	}

	return data, nil
}

// TocEntry is a file entry in the table of contents.
type TocEntry struct {
	Hash      Hash       // 64-bit FNV-1a hash of the file's path on disk
	Offset    uint32     // File's offset (in bytes) from the start of the nvc file
	RawLength uint32     // Length (in bytes) of the file after it has been extracted
	Length    uint32     // Length (in bytes) of file as it is stored in the nvc file
	Flags     EntryFlags // Indicates whether file is compressed or encrypted (TODO: determine if this is a bitmask)
}

func (e TocEntry) String() string {
	return fmt.Sprintf("%v offset=%v %vB (%vB on disk) flags=%v",
		e.Hash,
		e.Offset,
		e.RawLength,
		e.Length,
		e.Flags)
}

// readMagic reads the magic bytes at the beginning of a .nvc file.
func readMagic(r io.Reader) error {
	header := [8]byte{}
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return err
	}

	if string(header[:]) != magic {
		return ErrNoMagicFound
	}

	return nil
}

// readEntry reads a single ToC entry from r.
func readEntry(r io.Reader) (TocEntry, error) {
	var e TocEntry
	if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
		return e, err
	}
	return e, nil
}

// readCount reads the ToC entry count from r.
func readCount(r io.Reader) (uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return 0, err
	}
	return count, nil
}

// Hash is a 64-bit FNV-1a hash
type Hash uint64

// String2Hash returns the 64-bit FNV-1a hash of s
func String2Hash(s string) Hash {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return Hash(hash.Sum64())
}

func (h Hash) String() string {
	return fmt.Sprintf("%016x", uint64(h))
}

// Writer is an nvc archive writer.
type Writer struct {
	toc []TocEntry
	w   io.WriteSeeker

	// index keeps track of how many times Create has been called.
	// Since the table of contents is at the start of the archive,
	// and all archive member files are after the table of contents, the number of
	// files must be known ahead of time. Otherwise, writing the table of contents
	// would result in member file contents being partially overwritten.
	index int
}

// NewWriter returns an nvc archive writer that writes to w.
// length is the number of files that will be placed in the archive.
// Finalize should be called once all files have been written to the archive (via Create or CreateCompressed).
func NewWriter(w io.WriteSeeker, length uint32) (Writer, error) {
	// Start by writing 0s to w until the point at which the first file will start
	headerLen := 8 + 4 + (24 * length)
	_, err := w.Write(make([]byte, headerLen))
	if err != nil {
		return Writer{}, err
	}

	return Writer{
		toc:   make([]TocEntry, length),
		w:     w,
		index: 0,
	}, nil
}

// cumulativeWriter wraps an io.Writer and keeps a running total of how many bytes have been written
type cumulativeWriter struct {
	w     io.Writer
	count uint64
}

func (w *cumulativeWriter) Count() uint64 {
	return w.count
}

func (w *cumulativeWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.count += uint64(n)
	return n, err
}

// cumulativeReader wraps an io.Reader and keeps a running total of how many bytes have been read
type cumulativeReader struct {
	r     io.Reader
	count uint64
}

func (r *cumulativeReader) Count() uint64 {
	return r.count
}

func (r *cumulativeReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.count += uint64(n)
	return n, err
}

// Create reads an archive member file from r and writes it to w.
//
// Create increments w's internal Table of Contents entry counter by 1; it will panic if this counter exceeds the value of "length" that was passed to NewWriter.
// This function is not thread-safe; only one archive member file can be written to w at a time.
func (w *Writer) Create(r io.Reader, hash Hash) (int64, error) {
	if w.index == len(w.toc) {
		panic("File count exceeds originally specified number")
	}

	// Increment index at start of function so it won't get reused in the event of an early return
	idx := w.index
	w.index++

	reader := &cumulativeReader{r, 0}
	currentPos, _ := w.w.Seek(0, io.SeekCurrent)

	var entry *TocEntry = &(w.toc[idx])
	entry.Hash = hash
	entry.Offset = uint32(currentPos)
	entry.Flags = EntryFlagNoCompression

	written, err := io.Copy(w.w, reader)
	if err != nil {
		return written, err
	}
	read := reader.Count()

	entry.RawLength = uint32(read)
	entry.Length = uint32(written)

	return written, nil
}

// CreateCompressed reads an archive member file from r, compresses it using zlib compression, and writes it to w.
// See the documentation for [compress/zlib] for the acceptable values of level.

// CreateCompressed increments w's internal Table of Contents entry counter by 1; it will panic if this counter exceeds the value of "length" that was passed to NewWriter.
// This function is not thread-safe; only one archive member file can be written to w at a time.
func (w *Writer) CreateCompressed(r io.Reader, hash Hash, level int) (int64, error) {
	if w.index == len(w.toc) {
		panic("File count exceeds originally specified number")
	}

	// Increment index at start of function so it won't get reused in the event of an early return
	idx := w.index
	w.index++

	writer := cumulativeWriter{w.w, 0}
	zWriter, err := zlib.NewWriterLevel(&writer, level)
	if err != nil {
		return 0, err
	}

	reader := &cumulativeReader{r, 0}
	currentPos, _ := w.w.Seek(0, io.SeekCurrent)

	var entry *TocEntry = &(w.toc[idx])
	entry.Hash = hash
	entry.Offset = uint32(currentPos)
	entry.Flags = EntryFlagZlibCompression

	bytesWritten, err := io.Copy(zWriter, reader)
	if err != nil {
		return bytesWritten, err
	}

	zWriter.Close()

	bytesRead := reader.Count()
	bytesWritten = int64(writer.Count())

	entry.RawLength = uint32(bytesRead)
	entry.Length = uint32(bytesWritten)

	return bytesWritten, nil
}

// Finalize writes the nvc header to the start of w.
// It is an error to call Create after Finalize has been called.
func (w *Writer) Finalize() error {
	_, err := w.w.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = binary.Write(w.w, binary.LittleEndian, []byte(magic))
	if err != nil {
		return err
	}

	err = binary.Write(w.w, binary.LittleEndian, int32(len(w.toc)))
	if err != nil {
		return err
	}

	for _, entry := range w.toc {
		err = binary.Write(w.w, binary.LittleEndian, entry)
		if err != nil {
			return err
		}
	}

	return nil
}
