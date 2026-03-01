package egmanifest_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bur4ky/egmanifest"
)

//go:embed testdata/binary-valid.manifest
var binManifest []byte

func TestBinary(t *testing.T) {
	t.Parallel()

	bm, err := egmanifest.ParseBinary(binManifest)
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

func TestBinaryBadMagic(t *testing.T) {
	t.Parallel()

	bad := make([]byte, len(binManifest))
	copy(bad, binManifest)
	bad[0] = 0x00
	bad[1] = 0x00
	bad[2] = 0x00
	bad[3] = 0x00

	_, err := egmanifest.ParseBinary(bad)
	require.Error(t, err)
	require.Equal(t, egmanifest.ErrBadManifestMagic, err)
}

func BenchmarkBinaryParse(b *testing.B) {
	for b.Loop() {
		_, err := egmanifest.ParseBinary(binManifest)
		require.NoError(b, err)
	}
}
