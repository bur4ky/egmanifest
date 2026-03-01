package egmanifest

import (
	"encoding/binary"
	"fmt"
)

// encodeBlob encodes an integer value into a blob string of 3-digit decimal byte groups in little-endian order.
func encodeBlob[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64](v T) string {
	size := binary.Size(v)
	buf := make([]byte, size*3)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	for i := range size {
		buf[i*3] = '0' + b[i]/100
		buf[i*3+1] = '0' + (b[i]/10)%10
		buf[i*3+2] = '0' + b[i]%10
	}

	return string(buf)
}

// decodeBlob decodes a blob string of 3-digit decimal byte groups into an integer value in little-endian order.
func decodeBlob[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64](s string) (T, error) {
	size := binary.Size(T(0))
	if len(s) != size*3 {
		return 0, fmt.Errorf("invalid blob length: %d", len(s))
	}

	var raw uint64
	for i := range size {
		h := s[i*3 : i*3+3]
		v := int(h[0]-'0')*100 + int(h[1]-'0')*10 + int(h[2]-'0')
		if v > 255 {
			return 0, fmt.Errorf("invalid byte value: %d", v)
		}

		raw |= uint64(v) << (i * 8)
	}

	return T(raw), nil
}
