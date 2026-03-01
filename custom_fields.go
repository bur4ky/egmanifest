package egmanifest

import (
	"io"

	"github.com/bur4ky/egmanifest/binreader"
)

// CustomFields contains key-value metadata pairs from the manifest.
type CustomFields struct {
	DataSize    uint32
	DataVersion uint8
	Count       uint32
	Fields      map[string]string
}

// ReadCustomFields reads the custom fields, then seeks past any remaining
// bytes in the custom fields section based on the DataSize field.
func ReadCustomFields(reader *binreader.Reader) (*CustomFields, error) {
	var cf CustomFields
	var err error

	start := reader.Offset()
	cf.DataSize, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	cf.DataVersion, err = reader.Uint8()
	if err != nil {
		return nil, err
	}

	cf.Count, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	cf.Fields = make(map[string]string, cf.Count)

	keys := make([]string, cf.Count)
	for i := range cf.Count {
		keys[i], err = reader.FString()
		if err != nil {
			return nil, err
		}
	}

	for i := range cf.Count {
		val, err := reader.FString()
		if err != nil {
			return nil, err
		}

		cf.Fields[keys[i]] = val
	}

	_, err = reader.Seek(start+int64(cf.DataSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &cf, nil
}
