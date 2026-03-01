package egmanifest_test

import (
	"crypto/sha1"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bur4ky/egmanifest"
	"github.com/bur4ky/egmanifest/binreader"
)

//go:embed testdata/chunk-compressed.chunk
var chunkCompressed []byte

//go:embed testdata/chunk-uncompressed.chunk
var chunkPlain []byte

func TestChunkPlain(t *testing.T) {
	t.Parallel()

	_, err := egmanifest.ParseChunk(chunkPlain)
	require.NoError(t, err)
}

func TestChunkCompressed(t *testing.T) {
	t.Parallel()

	_, err := egmanifest.ParseChunk(chunkCompressed)
	require.NoError(t, err)
}

func TestChunkCompressedMatchesPlain(t *testing.T) {
	t.Parallel()

	reader := binreader.New(chunkCompressed)
	header, err := egmanifest.ReadChunkHeader(reader)
	require.NoError(t, err)

	data, err := egmanifest.ParseChunk(chunkCompressed)
	require.NoError(t, err)

	sum := sha1.Sum(data)
	require.Equal(t, header.SHAHash[:], sum[:])
}

func BenchmarkChunkPlain(b *testing.B) {
	for b.Loop() {
		_, err := egmanifest.ParseChunk(chunkPlain)
		require.NoError(b, err)
	}
}

func BenchmarkChunkCompressed(b *testing.B) {
	for b.Loop() {
		_, err := egmanifest.ParseChunk(chunkCompressed)
		require.NoError(b, err)
	}
}
