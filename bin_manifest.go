package egmanifest

import (
	"compress/zlib"
	"errors"
	"io"

	"github.com/bur4ky/egmanifest/binreader"
)

// BinaryManifestMagic is the expected magic number at the start of a binary manifest file.
const BinaryManifestMagic = 0x44BEC00C

var (
	// ErrBadManifestMagic is returned when the magic number at the start of the manifest file does not match the expected value.
	ErrBadManifestMagic  = errors.New("bad manifest magic found")
	ErrManifestEncrypted = errors.New("manifest is encrypted")
)

// BinaryManifest represents the entire manifest as read from a binary file.
type BinaryManifest struct {
	Header        *ManifestHeader
	Meta          *ManifestMeta
	ChunkDataList *ChunkDataList
	Files         *FileManifestList
	CustomFields  *CustomFields
}

// ParseBinary parses all sections of a manifest from the given byte slice, decompressing the data if needed.
func ParseBinary(b []byte) (*BinaryManifest, error) {
	var manifest BinaryManifest
	var err error

	reader := binreader.New(b)
	magic, err := reader.Uint32()
	if err != nil {
		return nil, err
	}

	if magic != BinaryManifestMagic {
		return nil, ErrBadManifestMagic
	}

	manifest.Header, err = ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	_, err = reader.Seek(int64(manifest.Header.HeaderSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	if manifest.Header.StoredAs&StoredAsCompressed != 0 {
		zr, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, manifest.Header.DataSizeUncompressed)
		_, err = io.ReadFull(zr, buf)
		_ = zr.Close()
		if err != nil {
			return nil, err
		}

		reader = binreader.New(buf)
	}

	if manifest.Header.StoredAs&StoredAsEncrypted != 0 {
		return nil, ErrManifestEncrypted
	}

	manifest.Meta, err = ReadMeta(reader)
	if err != nil {
		return nil, err
	}

	manifest.ChunkDataList, err = ReadChunkDataList(reader)
	if err != nil {
		return nil, err
	}

	manifest.Files, err = ReadFileManifestList(reader, manifest.ChunkDataList)
	if err != nil {
		return nil, err
	}

	manifest.CustomFields, err = ReadCustomFields(reader)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
