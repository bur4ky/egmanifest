package egmanifest

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/bur4ky/egmanifest/binreader"
)

// FileManifestList contains the list of all files in the manifest.
type FileManifestList struct {
	DataSize    uint32
	DataVersion uint8
	Count       uint32

	FileManifestList []FileManifest
}

// FileManifest describes a single file in the manifest.
type FileManifest struct {
	Filename      string
	SymlinkTarget string
	SHAHash       [sha1.Size]byte
	FileMetaFlags FileMetaFlag
	InstallTags   []string
	FileSize      uint32
	ChunkParts    []ChunkPart
	MimeType      string
}

// ChunkPart describes a contiguous region within a chunk that contributes to a file's content.
type ChunkPart struct {
	DataSize   uint32
	ParentGUID uuid.UUID
	Offset     uint32
	Size       uint32

	Chunk *Chunk
}

// ReadFileManifestList reads the file manifest list section,
// resolves chunk references against the provided ChunkDataList,
// then seeks past any remaining bytes based on the DataSize field.
func ReadFileManifestList(reader *binreader.Reader, dataList *ChunkDataList) (*FileManifestList, error) {
	var list FileManifestList
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

	list.FileManifestList = make([]FileManifest, list.Count)

	for i := range list.FileManifestList {
		list.FileManifestList[i].Filename, err = reader.FString()
		if err != nil {
			return nil, err
		}
	}

	for i := range list.FileManifestList {
		list.FileManifestList[i].SymlinkTarget, err = reader.FString()
		if err != nil {
			return nil, err
		}
	}

	for i := range list.FileManifestList {
		b, err := reader.Bytes(sha1.Size)
		if err != nil {
			return nil, err
		}

		list.FileManifestList[i].SHAHash = [sha1.Size]byte(b)
	}

	for i := range list.FileManifestList {
		flags, err := reader.Uint8()
		if err != nil {
			return nil, err
		}

		list.FileManifestList[i].FileMetaFlags = FileMetaFlag(flags)
	}

	for i := range list.FileManifestList {
		list.FileManifestList[i].InstallTags, err = reader.FStringArray()
		if err != nil {
			return nil, err
		}
	}

	for i := range list.FileManifestList {
		chunkPartsSize, err := reader.Uint32()
		if err != nil {
			return nil, err
		}

		fm := &list.FileManifestList[i]
		fm.ChunkParts = make([]ChunkPart, chunkPartsSize)

		var fileSize uint32
		for c := range fm.ChunkParts {
			cp := &fm.ChunkParts[c]

			cp.DataSize, err = reader.Uint32()
			if err != nil {
				return nil, err
			}

			cp.ParentGUID, err = reader.GUID()
			if err != nil {
				return nil, err
			}

			idx, ok := dataList.ChunkLookup[cp.ParentGUID]
			if !ok {
				return nil, fmt.Errorf("chunk GUID %s not found", cp.ParentGUID)
			}

			cp.Chunk = &dataList.Elements[idx]

			cp.Offset, err = reader.Uint32()
			if err != nil {
				return nil, err
			}

			cp.Size, err = reader.Uint32()
			if err != nil {
				return nil, err
			}

			fileSize += cp.Size
		}

		fm.FileSize = fileSize
	}

	if list.DataVersion >= 2 {
		for range list.Count {
			a, err := reader.Int32()
			if err != nil {
				return nil, err
			}

			_, err = reader.Seek(int64(a)*16, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}

		for i := range list.Count {
			list.FileManifestList[i].MimeType, err = reader.FString()
			if err != nil {
				return nil, err
			}
		}

		for range list.Count {
			_, err = reader.Seek(32, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err = reader.Seek(start+int64(list.DataSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &list, nil
}
