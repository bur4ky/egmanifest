package egmanifest

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

// JSONManifest is the JSON representation of a manifest file, containing app metadata and the list of file entries.
type JSONManifest struct {
	ManifestFileVersion string                 `json:"ManifestFileVersion"`
	IsFileData          *bool                  `json:"bIsFileData,omitempty"`
	AppID               uint32                 `json:"AppID"`
	AppNameString       string                 `json:"AppNameString"`
	BuildVersionString  string                 `json:"BuildVersionString"`
	LaunchExeString     string                 `json:"LaunchExeString"`
	LaunchCommand       string                 `json:"LaunchCommand"`
	PrereqIDs           []string               `json:"PrereqIds,omitempty"`
	PrereqName          string                 `json:"PrereqName"`
	PrereqPath          string                 `json:"PrereqPath"`
	PrereqArgs          string                 `json:"PrereqArgs"`
	FileManifestList    []JSONFileManifestList `json:"FileManifestList"`
	ChunkHashList       map[string]string      `json:"ChunkHashList,omitempty"`
	ChunkShaList        map[string]string      `json:"ChunkShaList,omitempty"`
	DataGroupList       map[string]string      `json:"DataGroupList,omitempty"`
	ChunkFilesizeList   map[string]string      `json:"ChunkFilesizeList,omitempty"`
	CustomFields        map[string]string      `json:"CustomFields,omitempty"`
}

// ToBinaryManifest converts the JSON manifest into a new BinaryManifest.
func (m JSONManifest) ToBinaryManifest() (*BinaryManifest, error) {
	header := &ManifestHeader{
		StoredAs: StoredAsPlainText,
	}

	if m.ManifestFileVersion == "" {
		header.Version = FeatureLevelCustomFields
	} else {
		version, err := decodeBlob[uint32](m.ManifestFileVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid ManifestFileVersion: %w", err)
		}

		header.Version = FeatureLevel(int32(version))
		if header.Version == FeatureLevelBrokenJSONVersion {
			header.Version = FeatureLevelStoresChunkFileSizes
		}
	}

	meta := &ManifestMeta{
		AppID:         int32(m.AppID),
		AppName:       m.AppNameString,
		BuildVersion:  m.BuildVersionString,
		LaunchExe:     m.LaunchExeString,
		LaunchCommand: m.LaunchCommand,
		PrereqIDs:     m.PrereqIDs,
		PrereqName:    m.PrereqName,
		PrereqPath:    m.PrereqPath,
		PrereqArgs:    m.PrereqArgs,
	}

	chunks := make(map[uuid.UUID]*Chunk)
	for _, file := range m.FileManifestList {
		for _, cp := range file.FileChunkParts {
			if _, exists := chunks[cp.GUID]; exists {
				continue
			}

			sum := sha1.Sum([]byte(cp.GUID.String()))
			chunks[cp.GUID] = &Chunk{
				GUID:       cp.GUID,
				Hash:       binary.LittleEndian.Uint64(sum[:8]),
				SHAHash:    sum,
				WindowSize: LegacyFixedChunkWindow,
				FileSize:   LegacyFixedChunkWindow,
			}
		}
	}

	hasChunkHashList := false
	for guidStr, hashStr := range m.ChunkHashList {
		guid, err := uuid.Parse(guidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid GUID in ChunkHashList: %w", err)
		}

		hash, err := decodeBlob[uint64](hashStr)
		if err != nil {
			return nil, fmt.Errorf("invalid hash in ChunkHashList: %w", err)
		}

		if c, ok := chunks[guid]; ok {
			c.Hash = hash
		}

		hasChunkHashList = true
	}

	for guidStr, shaStr := range m.ChunkShaList {
		guid, err := uuid.Parse(guidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid GUID in ChunkShaList: %w", err)
		}

		b, err := hex.DecodeString(shaStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SHA in ChunkShaList: %w", err)
		}

		if len(b) != sha1.Size {
			return nil, fmt.Errorf("invalid SHA length in ChunkShaList: %d", len(b))
		}

		if c, ok := chunks[guid]; ok {
			c.SHAHash = [sha1.Size]byte(b)
		}
	}

	for guidStr, groupStr := range m.DataGroupList {
		guid, err := uuid.Parse(guidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid GUID in DataGroupList: %w", err)
		}

		group, err := decodeBlob[uint8](groupStr)
		if err != nil {
			return nil, fmt.Errorf("invalid group in DataGroupList: %w", err)
		}

		if c, ok := chunks[guid]; ok {
			c.Group = group
		}
	}

	if len(m.ChunkFilesizeList) > 0 {
		for guidStr, sizeStr := range m.ChunkFilesizeList {
			guid, err := uuid.Parse(guidStr)
			if err != nil {
				return nil, fmt.Errorf("invalid GUID in ChunkFilesizeList: %w", err)
			}

			fileSize, err := decodeBlob[int64](sizeStr)
			if err != nil {
				return nil, fmt.Errorf("invalid file size in ChunkFilesizeList: %w", err)
			}

			if c, ok := chunks[guid]; ok {
				c.FileSize = uint64(fileSize)
			}
		}
	}

	if m.IsFileData != nil {
		meta.IsFileData = *m.IsFileData
	} else {
		meta.IsFileData = !hasChunkHashList
	}

	var customFields *CustomFields
	if len(m.CustomFields) > 0 {
		customFields = &CustomFields{
			Count:  uint32(len(m.CustomFields)),
			Fields: m.CustomFields,
		}
	}

	chunkLookup := make(map[uuid.UUID]uint32, len(chunks))
	chunkElements := make([]Chunk, len(chunks))
	idx := uint32(0)
	for id, chunk := range chunks {
		chunkLookup[id] = idx
		chunkElements[idx] = *chunk
		idx++
	}

	chunkDataList := &ChunkDataList{
		Count:       uint32(len(chunks)),
		Elements:    chunkElements,
		ChunkLookup: chunkLookup,
	}

	files := make([]FileManifest, 0, len(m.FileManifestList))
	for _, fm := range m.FileManifestList {
		shaHash, err := parseFileHash(fm.FileHash)
		if err != nil {
			return nil, err
		}

		var flags FileMetaFlag
		switch {
		case fm.IsUnixExecutable:
			flags |= FileMetaFlagUnixExecutable
		case fm.IsReadOnly:
			flags |= FileMetaFlagReadOnly
		case fm.IsCompressed:
			flags |= FileMetaFlagCompressed
		}

		chunkParts := make([]ChunkPart, len(fm.FileChunkParts))
		fileSize := uint32(0)

		for j, cp := range fm.FileChunkParts {
			fileSize += cp.Size
			chunkParts[j] = ChunkPart{
				ParentGUID: cp.GUID,
				Offset:     cp.Offset,
				Size:       cp.Size,
				Chunk:      chunks[cp.GUID],
			}
		}

		files = append(files, FileManifest{
			Filename:      fm.Filename,
			SymlinkTarget: fm.SymlinkTarget,
			InstallTags:   fm.InstallTags,
			SHAHash:       shaHash,
			FileMetaFlags: flags,
			ChunkParts:    chunkParts,
			FileSize:      fileSize,
		})
	}

	fileManifestList := &FileManifestList{
		Count:            uint32(len(files)),
		FileManifestList: files,
	}

	return &BinaryManifest{
		Header:        header,
		Meta:          meta,
		ChunkDataList: chunkDataList,
		Files:         fileManifestList,
		CustomFields:  customFields,
	}, nil
}

// MarshalJSON encodes the manifest, serializing AppID as a decimal string per the manifest JSON format.
func (m JSONManifest) MarshalJSON() ([]byte, error) {
	type Alias JSONManifest
	return json.Marshal(&struct {
		AppID string `json:"AppID"`
		*Alias
	}{
		AppID: encodeBlob(m.AppID),
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON decodes the manifest, parsing AppID from a decimal string per the manifest JSON format.
func (m *JSONManifest) UnmarshalJSON(b []byte) error {
	type Alias JSONManifest

	aux := &struct {
		AppID string `json:"AppID"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	err := json.Unmarshal(b, aux)
	if err != nil {
		return err
	}

	m.AppID, err = decodeBlob[uint32](aux.AppID)
	if err != nil {
		return err
	}

	return nil
}

// ParseJSON parses a JSON manifest from the given byte slice.
func ParseJSON(data []byte) (*JSONManifest, error) {
	var manifest JSONManifest
	err := json.Unmarshal(data, &manifest)
	return &manifest, err
}

func parseFileHash(s string) ([sha1.Size]byte, error) {
	var out [sha1.Size]byte
	if len(s) != sha1.Size*3 {
		return out, fmt.Errorf("invalid file hash length: %d", len(s))
	}

	for i := range sha1.Size {
		start := i * 3
		b, err := strconv.ParseUint(s[start:start+3], 10, 8)
		if err != nil {
			return out, fmt.Errorf("invalid file hash byte: %w", err)
		}

		out[i] = byte(b)
	}

	return out, nil
}
