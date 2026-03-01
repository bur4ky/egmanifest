package egmanifest

import (
	"crypto/sha1"
	"errors"
	"io"

	"github.com/google/uuid"

	"github.com/bur4ky/egmanifest/binreader"
)

// ChunkHeaderMagic is the expected magic number at the start of a chunk header.
const ChunkHeaderMagic = 0xB1FE3AA2

// ErrBadChunkHeaderMagic is returned when the magic number at the start of a chunk header does not match the expected value.
var ErrBadChunkHeaderMagic = errors.New("bad chunk header magic found")

// ChunkHeader contains metadata for a chunk file.
type ChunkHeader struct {
	Version              ChunkHeaderVersion
	HeaderSize           uint32
	DataSizeCompressed   uint32
	DataSizeUncompressed uint32
	GUID                 uuid.UUID
	RollingHash          uint64
	StoredAs             StoredAs
	SHAHash              [sha1.Size]byte
	HashType             uint32
}

// ReadChunkHeader reads and parses a chunk file header, then seeks past
// any remaining bytes in the header section based on the HeaderSize field.
func ReadChunkHeader(reader *binreader.Reader) (*ChunkHeader, error) {
	var header ChunkHeader

	start := reader.Offset()
	magic, err := reader.Uint32()
	if err != nil {
		return nil, err
	}

	if magic != ChunkHeaderMagic {
		return nil, ErrBadChunkHeaderMagic
	}

	version, err := reader.Uint32()
	if err != nil {
		return nil, err
	}

	header.Version = ChunkHeaderVersion(version)

	header.HeaderSize, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	header.DataSizeCompressed, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	header.GUID, err = reader.GUID()
	if err != nil {
		return nil, err
	}

	header.RollingHash, err = reader.Uint64()
	if err != nil {
		return nil, err
	}

	storedAs, err := reader.Uint8()
	if err != nil {
		return nil, err
	}

	header.StoredAs = StoredAs(storedAs)

	if header.Version >= ChunkHeaderVersionStoresShaAndHashType {
		sha, err := reader.Bytes(sha1.Size)
		if err != nil {
			return nil, err
		}

		header.SHAHash = [sha1.Size]byte(sha)

		header.HashType, err = reader.Uint32()
		if err != nil {
			return nil, err
		}
	}

	if header.Version >= ChunkHeaderVersionStoresDataSizeUncompressed {
		header.DataSizeUncompressed, err = reader.Uint32()
		if err != nil {
			return nil, err
		}
	} else {
		header.DataSizeUncompressed = LegacyFixedChunkWindow
	}

	_, err = reader.Seek(start+int64(header.HeaderSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &header, err
}
