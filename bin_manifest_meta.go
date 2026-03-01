package egmanifest

import (
	"io"

	"github.com/bur4ky/egmanifest/binreader"
)

// ManifestMeta represents manifets metadata.
type ManifestMeta struct {
	DataSize    uint32
	DataVersion ManifestMetaVersion

	FeatureLevel  FeatureLevel
	IsFileData    bool
	AppID         int32
	AppName       string
	BuildVersion  string
	LaunchExe     string
	LaunchCommand string
	PrereqIDs     []string
	PrereqName    string
	PrereqPath    string
	PrereqArgs    string

	BuildID             string
	UninstallActionPath string
	UninstallActionArgs string
}

// ReadMeta reads the manifest metadata, then seeks past any remaining
// bytes in the metadata section based on the DataSize field.
func ReadMeta(reader *binreader.Reader) (*ManifestMeta, error) {
	var meta ManifestMeta
	var err error

	start := reader.Offset()
	meta.DataSize, err = reader.Uint32()
	if err != nil {
		return nil, err
	}

	dataVersion, err := reader.Uint8()
	if err != nil {
		return nil, err
	}

	meta.DataVersion = ManifestMetaVersion(dataVersion)

	featureLevel, err := reader.Int32()
	if err != nil {
		return nil, err
	}

	meta.FeatureLevel = FeatureLevel(featureLevel)

	meta.IsFileData, err = reader.Bool()
	if err != nil {
		return nil, err
	}

	meta.AppID, err = reader.Int32()
	if err != nil {
		return nil, err
	}

	meta.AppName, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.BuildVersion, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.LaunchExe, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.LaunchCommand, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.PrereqIDs, err = reader.FStringArray()
	if err != nil {
		return nil, err
	}

	meta.PrereqName, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.PrereqPath, err = reader.FString()
	if err != nil {
		return nil, err
	}

	meta.PrereqArgs, err = reader.FString()
	if err != nil {
		return nil, err
	}

	if meta.DataVersion >= ManifestMetaVersionSerialisesBuildID {
		meta.BuildID, err = reader.FString()
		if err != nil {
			return nil, err
		}
	}

	if meta.DataVersion > ManifestMetaVersionSerialisesBuildID {
		meta.UninstallActionPath, err = reader.FString()
		if err != nil {
			return nil, err
		}

		meta.UninstallActionArgs, err = reader.FString()
		if err != nil {
			return nil, err
		}
	}

	_, err = reader.Seek(start+int64(meta.DataSize), io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
