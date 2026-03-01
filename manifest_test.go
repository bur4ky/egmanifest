package egmanifest_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bur4ky/egmanifest"
)

func TestManifestDetectJSON(t *testing.T) {
	t.Parallel()

	m, err := egmanifest.Parse(jsonManifest)
	require.NoError(t, err)
	require.NotNil(t, m)

	require.NotNil(t, m.Header)
	require.NotNil(t, m.Meta)
	require.NotNil(t, m.ChunkDataList)
	require.NotNil(t, m.Files)
}

func TestManifestDetectBinary(t *testing.T) {
	t.Parallel()

	m, err := egmanifest.Parse(binManifest)
	require.NoError(t, err)
	require.NotNil(t, m)

	require.NotNil(t, m.Header)
	require.NotNil(t, m.Meta)
	require.NotNil(t, m.ChunkDataList)
	require.NotNil(t, m.Files)
}

func TestManifestInvalidJSON(t *testing.T) {
	t.Parallel()

	bad := []byte(`{invalid json}`)
	m, err := egmanifest.Parse(bad)
	require.Error(t, err)
	require.Nil(t, m)
}

func TestManifestInvalidBinary(t *testing.T) {
	t.Parallel()

	bad := []byte{0x42, 0x11, 0x39, 0xFF}
	m, err := egmanifest.Parse(bad)
	require.Error(t, err)
	require.Nil(t, m)
}

func TestManifestEmptyInput(t *testing.T) {
	t.Parallel()

	m, err := egmanifest.Parse([]byte{})
	require.Error(t, err)
	require.Nil(t, m)
}
