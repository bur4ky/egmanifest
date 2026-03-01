package binreader

import (
	"encoding/binary"
	"errors"
	"unicode/utf16"

	"github.com/google/uuid"
)

// maxStringBytes is an arbitrary limit to prevent OOM errors when reading malformed strings.
const maxStringBytes = 1 << 30

// GUID reads a 16-byte GUID and returns it as a uuid.UUID.
func (r *Reader) GUID() (uuid.UUID, error) {
	b, err := r.Bytes(16)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	for i := range 4 {
		v := binary.BigEndian.Uint32(b[i*4:])
		binary.LittleEndian.PutUint32(id[i*4:], v)
	}

	return id, nil
}

// FString reads an FString.
// Positive length indicates UTF-8 (with null terminator), negative length indicates UTF-16LE.
func (r *Reader) FString() (string, error) {
	n, err := r.Int32()
	if err != nil {
		return "", err
	}

	switch {
	case n == 0:
		return "", nil
	case n > 0:
		size := int(n)
		if size > maxStringBytes {
			return "", errors.New("string too large")
		}

		buf, err := r.Bytes(size)
		if err != nil {
			return "", err
		}

		if buf[size-1] == 0 {
			buf = buf[:size-1]
		}

		return string(buf), nil
	default:
		chars := -int(n)
		byteLen := chars * 2
		if byteLen > maxStringBytes {
			return "", errors.New("string too large")
		}

		buf, err := r.Bytes(byteLen)
		if err != nil {
			return "", err
		}

		u16 := make([]uint16, chars)
		for i := range chars {
			u16[i] = binary.LittleEndian.Uint16(buf[i*2:])
		}

		if u16[chars-1] == 0 {
			u16 = u16[:chars-1]
		}

		return string(utf16.Decode(u16)), nil
	}
}

// FStringArray reads an array of FStrings.
func (r *Reader) FStringArray() ([]string, error) {
	n, err := r.Uint32()
	if err != nil {
		return nil, err
	}

	out := make([]string, n)
	for i := range out {
		out[i], err = r.FString()
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
