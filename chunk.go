package egmanifest

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/bur4ky/egmanifest/binreader"
)

// Chunk represents a single chunk of data as defined in the manifest's ChunkDataList.
type Chunk struct {
	GUID       uuid.UUID
	Hash       uint64
	SHAHash    [sha1.Size]byte
	Group      uint8
	WindowSize uint32
	FileSize   uint64
}

// URL generates the full download URL for the chunk.
//
// The chunksDirURL parameter should be the base chunks directory URL,
// including the version-specific subdirectory.
// For example: http://epicgames-download1.akamaized.net/Builds/Fortnite/CloudDir/ChunksV4
//
// To get the correct chunk subdirectory (e.g., "ChunksV4"),
// use BinaryManifest.Header.Version.ChunkSubDir().
func (c *Chunk) URL(chunksDirURL string) string {
	return fmt.Sprintf("%s/%02d/%016X_%X.chunk", chunksDirURL, c.Group, c.Hash, c.GUID[:])
}

// ParseChunk parses a chunk file, reading the header and data sections
// and decompressing the data if needed.
func ParseChunk(data []byte) ([]byte, error) {
	reader := binreader.New(data)
	header, err := ReadChunkHeader(reader)
	if err != nil {
		return nil, err
	}

	b, err := reader.Bytes(int(header.DataSizeCompressed))
	if err != nil {
		return nil, err
	}

	switch header.StoredAs {
	case StoredAsPlainText:
		return b, nil
	case StoredAsCompressed:
		return DecompressChunk(b)
	case StoredAsEncrypted:
		return nil, fmt.Errorf("chunk is encrypted")
	default:
		return nil, fmt.Errorf("unknown storage mode %d", header.StoredAs)
	}
}

// DecompressChunk takes compressed chunk data and returns the decompressed data.
func DecompressChunk(data []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer zr.Close()
	return io.ReadAll(zr)
}
