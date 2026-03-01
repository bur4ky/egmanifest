// Package binreader provides a binary reader for reading various data types from a byte slice.
package binreader

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// Reader is a binary reader that reads from a byte slice. It implements io.ReadSeeker.
type Reader struct {
	data []byte
	off  int64
}

// New creates a new Reader for the given byte slice.
func New(data []byte) *Reader {
	return &Reader{data: data}
}

// Offset returns the current offset of the reader.
func (r *Reader) Offset() int64 {
	return r.off
}

// Read reads up to len(p) bytes into p.
func (r *Reader) Read(p []byte) (int, error) {
	if r.off >= r.len() {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.off:])
	r.off += int64(n)
	return n, nil
}

// Seek sets the offset for the next Read to offset, interpreted according to whence.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		// unchanged
	case io.SeekCurrent:
		offset = r.off + offset
	case io.SeekEnd:
		offset = r.len() + offset
	default:
		return 0, errors.New("invalid whence")
	}

	if offset < 0 || offset > r.len() {
		return 0, io.ErrUnexpectedEOF
	}

	r.off = offset
	return offset, nil
}

// Peek reads n bytes without advancing the offset.
func (r *Reader) Peek(n int) ([]byte, error) {
	offset := r.off
	b, err := r.Bytes(n)
	if err != nil {
		return nil, err
	}

	r.off = offset
	return b, nil
}

// Bytes reads n bytes and advances the offset.
func (r *Reader) Bytes(n int) ([]byte, error) {
	if n < 0 || int64(n) > (r.len()-r.off) {
		return nil, io.ErrUnexpectedEOF
	}

	start := r.off
	r.off += int64(n)
	return r.data[start:r.off], nil
}

// Bool reads a boolean value (1 byte) and returns true if it's non-zero.
func (r *Reader) Bool() (bool, error) {
	v, err := r.Uint8()
	return v != 0, err
}

// Uint8 reads an unsigned 8-bit integer.
func (r *Reader) Uint8() (uint8, error) {
	b, err := r.Bytes(1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

// Uint16 reads an unsigned 16-bit integer.
func (r *Reader) Uint16() (uint16, error) {
	b, err := r.Bytes(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(b), nil
}

// Uint32 reads an unsigned 32-bit integer.
func (r *Reader) Uint32() (uint32, error) {
	b, err := r.Bytes(4)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

// Uint64 reads an unsigned 64-bit integer.
func (r *Reader) Uint64() (uint64, error) {
	b, err := r.Bytes(8)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(b), nil
}

// Int8 reads an unsigned 8-bit integer and returns it as a signed 8-bit integer.
func (r *Reader) Int8() (int8, error) {
	v, err := r.Uint8()
	return int8(v), err
}

// Int16 reads an unsigned 16-bit integer and returns it as a signed 16-bit integer.
func (r *Reader) Int16() (int16, error) {
	v, err := r.Uint16()
	return int16(v), err
}

// Int32 reads an unsigned 32-bit integer and returns it as a signed 32-bit integer.
func (r *Reader) Int32() (int32, error) {
	v, err := r.Uint32()
	return int32(v), err
}

// Int64 reads an unsigned 64-bit integer and returns it as a signed 64-bit integer.
func (r *Reader) Int64() (int64, error) {
	v, err := r.Uint64()
	return int64(v), err
}

// Float32 reads a 32-bit float by reading its bits as a uint32 and converting it to float32.
func (r *Reader) Float32() (float32, error) {
	v, err := r.Uint32()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(v), nil
}

// Float64 reads a 64-bit float by reading its bits as a uint64 and converting it to float64.
func (r *Reader) Float64() (float64, error) {
	v, err := r.Uint64()
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(v), nil
}

// len returns the total length of the data in the reader.
func (r *Reader) len() int64 {
	return int64(len(r.data))
}
