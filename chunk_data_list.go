package egmanifest

import (
	"crypto/sha1"
	"io"

	"github.com/google/uuid"

	"github.com/bur4ky/egmanifest/binreader"
)

// ChunkDataList represents the list of chunks defined in the manifest.
type ChunkDataList struct {
	DataSize    uint32
	DataVersion uint8
	Count       uint32

	Elements    []Chunk
	ChunkLookup map[uuid.UUID]uint32
}

// ReadChunkDataList reads the chunk data list, then seeks past any remaining
// bytes in the chunk data section based on the DataSize field.
func ReadChunkDataList(reader *binreader.Reader) (*ChunkDataList, error) {
	var list ChunkDataList
	var err error

	start := reader.Offset()
	list.DataSize, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	list.DataVersion, err = reader.Uint8()
	if err != nil {
		return nil, err
	}

	list.Count, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	list.Elements = make([]Chunk, list.Count)
	list.ChunkLookup = make(map[uuid.UUID]uint32, list.Count)

	for i := range list.Count {
		guid, err := reader.GUID()
		if err != nil {
			return nil, err
		}

		list.Elements[i].GUID = guid
		list.ChunkLookup[guid] = i
	}

	for i := range list.Count {
		hash, err := reader.Uint64()
		if err != nil {
			return nil, err
		}

		list.Elements[i].Hash = hash
	}

	for i := range list.Count {
		sha, err := reader.Bytes(sha1.Size)
		if err != nil {
			return nil, err
		}

		list.Elements[i].SHAHash = [sha1.Size]byte(sha)
	}

	for i := range list.Count {
		group, err := reader.Uint8()
		if err != nil {
			return nil, err
		}

		list.Elements[i].Group = group
	}

	for i := range list.Count {
		windowSize, err := reader.Uint32()
		if err != nil {
			return nil, err
		}

		list.Elements[i].WindowSize = windowSize
	}

	for i := range list.Count {
		fileSize, err := reader.Uint64()
		if err != nil {
			return nil, err
		}

		list.Elements[i].FileSize = fileSize
	}

	_, err = reader.Seek(start+int64(list.DataSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &list, nil
}
