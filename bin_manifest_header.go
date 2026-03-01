package egmanifest

import (
	"crypto/sha1"
	"io"

	"github.com/bur4ky/egmanifest/binreader"
)

// ManifestHeader represents the header of a manifest file.
type ManifestHeader struct {
	HeaderSize           int32
	DataSizeUncompressed int32
	DataSizeCompressed   int32
	SHAHash              [sha1.Size]byte
	StoredAs             StoredAs
	Version              FeatureLevel
}

// ReadHeader reads the manifest header, then seeks past any remaining
// bytes in the header section based on the HeaderSize field.
func ReadHeader(reader *binreader.Reader) (*ManifestHeader, error) {
	var header ManifestHeader
	var err error

	start := reader.Offset()
	header.HeaderSize, err = reader.Int32()
	if err != nil {
		return nil, err
	}

	header.DataSizeUncompressed, err = reader.Int32()
	if err != nil {
		return nil, err
	}

	header.DataSizeCompressed, err = reader.Int32()
	if err != nil {
		return nil, err
	}

	sha, err := reader.Bytes(sha1.Size)
	if err != nil {
		return nil, err
	}

	header.SHAHash = [sha1.Size]byte(sha)

	storedAs, err := reader.Uint8()
	if err != nil {
		return nil, err
	}

	header.StoredAs = StoredAs(storedAs)

	if header.HeaderSize > ManifestHeaderVersionSizes[FeatureLevelOriginal] {
		v, err := reader.Int32()
		if err != nil {
			return nil, err
		}

		header.Version = FeatureLevel(v)
	} else {
		header.Version = FeatureLevelStoredAsCompressedUClass
	}

	_, err = reader.Seek(start+int64(header.HeaderSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &header, nil
}
