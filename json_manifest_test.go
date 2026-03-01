package egmanifest_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bur4ky/egmanifest"
)

//go:embed testdata/json-valid.manifest
var jsonManifest []byte

func TestJSONToBinary(t *testing.T) {
	t.Parallel()

	jm, err := egmanifest.ParseJSON(jsonManifest)
	require.NoError(t, err)

	bm, err := jm.ToBinaryManifest()
	require.NoError(t, err)
	require.NotNil(t, bm)

	require.NotNil(t, bm.Header)
	require.NotNil(t, bm.Meta)
	require.NotNil(t, bm.ChunkDataList)
	require.NotNil(t, bm.Files)

	require.NotZero(t, bm.ChunkDataList.Count)
	require.NotZero(t, bm.Files.Count)

	require.Equal(t, uint32(len(bm.ChunkDataList.Elements)), bm.ChunkDataList.Count)
	require.Equal(t, uint32(len(bm.Files.FileManifestList)), bm.Files.Count)
}

func BenchmarkJSONParse(b *testing.B) {
	for b.Loop() {
		_, err := egmanifest.ParseJSON(jsonManifest)
		require.NoError(b, err)
	}
}

func BenchmarkJSONToBinary(b *testing.B) {
	jm, err := egmanifest.ParseJSON(jsonManifest)
	require.NoError(b, err)

	b.ResetTimer()

	for b.Loop() {
		_, err = jm.ToBinaryManifest()
		require.NoError(b, err)
	}
}
