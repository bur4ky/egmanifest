package egmanifest

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

// JSONFileManifestList is the JSON representation of a file manifest entry.
type JSONFileManifestList struct {
	Filename         string          `json:"Filename"`
	FileHash         string          `json:"FileHash"`
	IsUnixExecutable bool            `json:"bIsUnixExecutable,omitempty"`
	IsReadOnly       bool            `json:"bIsReadOnly,omitempty"`
	IsCompressed     bool            `json:"bIsCompressed,omitempty"`
	SymlinkTarget    string          `json:"SymlinkTarget,omitempty"`
	InstallTags      []string        `json:"InstallTags,omitempty"`
	FileChunkParts   []JSONChunkPart `json:"FileChunkParts"`
}

// JSONChunkPart is the JSON representation of a chunk part, referencing a chunk by GUID with an offset and size.
type JSONChunkPart struct {
	GUID   uuid.UUID `json:"Guid"`
	Offset uint32    `json:"Offset"`
	Size   uint32    `json:"Size"`
}

// MarshalJSON encodes the chunk part fields as hex/decimal strings per the manifest JSON format.
func (c JSONChunkPart) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		GUID   string `json:"Guid"`
		Offset string `json:"Offset"`
		Size   string `json:"Size"`
	}{
		GUID:   strings.ToUpper(strings.ReplaceAll(c.GUID.String(), "-", "")),
		Offset: encodeBlob(c.Offset),
		Size:   encodeBlob(c.Size),
	})
}

// UnmarshalJSON decodes the chunk part fields from hex/decimal strings per the manifest JSON format.
func (c *JSONChunkPart) UnmarshalJSON(b []byte) error {
	aux := &struct {
		GUID   string `json:"Guid"`
		Offset string `json:"Offset"`
		Size   string `json:"Size"`
	}{}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	c.GUID, err = uuid.Parse(aux.GUID)
	if err != nil {
		return err
	}

	c.Offset, err = decodeBlob[uint32](aux.Offset)
	if err != nil {
		return err
	}

	c.Size, err = decodeBlob[uint32](aux.Size)
	if err != nil {
		return err
	}

	return nil
}
