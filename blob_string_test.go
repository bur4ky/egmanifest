package egmanifest

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlobUint32(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    uint32
		expected string
	}{
		{0, "000000000000"},
		{1 << 0, "001000000000"},
		{1 << 8, "000001000000"},
		{1 << 16, "000000001000"},
		{1 << 24, "000000000001"},
		{math.MaxUint32, "255255255255"},
		{518150, "006232007000"},
	}
	for _, c := range cases {
		t.Run(strconv.FormatUint(uint64(c.input), 10), func(t *testing.T) {
			t.Parallel()

			encoded := encodeBlob(c.input)
			require.Equal(t, c.expected, encoded)

			decoded, err := decodeBlob[uint32](encoded)
			require.NoError(t, err)

			require.Equal(t, c.input, decoded)
		})
	}
}

func TestBlobUint64(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    uint64
		expected string
	}{
		{0, "000000000000000000000000"},
		{1 << 0, "001000000000000000000000"},
		{1 << 8, "000001000000000000000000"},
		{1 << 32, "000000000000001000000000"},
		{math.MaxUint64, "255255255255255255255255"},
	}
	for _, c := range cases {
		t.Run(strconv.FormatUint(c.input, 10), func(t *testing.T) {
			t.Parallel()

			encoded := encodeBlob(c.input)
			require.Equal(t, c.expected, encoded)

			decoded, err := decodeBlob[uint64](encoded)
			require.NoError(t, err)

			require.Equal(t, c.input, decoded)
		})
	}
}

func TestBlobUint8(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    uint8
		expected string
	}{
		{0, "000"},
		{1, "001"},
		{255, "255"},
		{128, "128"},
	}
	for _, c := range cases {
		t.Run(strconv.FormatUint(uint64(c.input), 10), func(t *testing.T) {
			t.Parallel()

			encoded := encodeBlob(c.input)
			require.Equal(t, c.expected, encoded)

			decoded, err := decodeBlob[uint8](encoded)
			require.NoError(t, err)

			require.Equal(t, c.input, decoded)
		})
	}
}

func TestBlobInt64(t *testing.T) {
	t.Parallel()
	cases := []struct {
		input    int64
		expected string
	}{
		{0, "000000000000000000000000"},
		{1, "001000000000000000000000"},
		{-1, "255255255255255255255255"},
		{math.MinInt64, "000000000000000000000128"},
		{math.MaxInt64, "255255255255255255255127"},
		{1048576, "000000016000000000000000"},
	}
	for _, c := range cases {
		t.Run(strconv.FormatInt(c.input, 10), func(t *testing.T) {
			t.Parallel()

			encoded := encodeBlob(c.input)
			require.Equal(t, c.expected, encoded)

			decoded, err := decodeBlob[int64](encoded)
			require.NoError(t, err)

			require.Equal(t, c.input, decoded)
		})
	}
}

func TestDecodeBlobErrors(t *testing.T) {
	t.Parallel()

	t.Run("wrong length", func(t *testing.T) {
		t.Parallel()
		_, err := decodeBlob[uint32]("000000000")
		require.Error(t, err)
	})

	t.Run("invalid digit", func(t *testing.T) {
		t.Parallel()
		_, err := decodeBlob[uint32]("999000000000")
		require.Error(t, err)
	})

	t.Run("value over 255", func(t *testing.T) {
		t.Parallel()
		_, err := decodeBlob[uint32]("256000000000")
		require.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		_, err := decodeBlob[uint32]("")
		require.Error(t, err)
	})
}
