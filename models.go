package egmanifest

// LegacyFixedChunkWindow is the legacy fixed chunk window size, which was 1MiB.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ChunkData.h#L15
const LegacyFixedChunkWindow = 1024 * 1024

// FeatureLevel represents the feature level of a manifest.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Public/BuildPatchFeatureLevel.h#L12-L71
type FeatureLevel int32

const (
	// FeatureLevelOriginal is the first version of the format.
	FeatureLevelOriginal FeatureLevel = iota

	// FeatureLevelCustomFields is the version that adds support for custom fields.
	FeatureLevelCustomFields

	// FeatureLevelStartStoringVersion is the point where the manifest began storing its own version number.
	FeatureLevelStartStoringVersion

	// FeatureLevelDataFileRenames is when data files started including a hash in their filenames, sending these chunks
	// to ChunksV2.
	FeatureLevelDataFileRenames

	// FeatureLevelStoresIfChunkOrFileData is the version where the manifest records whether a build uses chunk-based
	// or file-based data.
	FeatureLevelStoresIfChunkOrFileData

	// FeatureLevelStoresDataGroupNumbers is the version where each chunk or file entry includes its data group number
	// for external readers.
	FeatureLevelStoresDataGroupNumbers

	// FeatureLevelChunkCompressionSupport is when chunk compression was added, storing these chunks in ChunksV3.
	// FileManifest data compression wasn’t added yet.
	FeatureLevelChunkCompressionSupport

	// FeatureLevelStoresPrerequisitesInfo is when the manifest began storing product prerequisite information.
	FeatureLevelStoresPrerequisitesInfo

	// FeatureLevelStoresChunkFileSizes is the version where chunk download sizes started being recorded.
	FeatureLevelStoresChunkFileSizes

	// FeatureLevelStoredAsCompressedUClass is when the manifest gained optional UObject-based serialization with
	// compression.
	FeatureLevelStoredAsCompressedUClass

	// FeatureLevelUnused0 was removed and never used.
	FeatureLevelUnused0

	// FeatureLevelUnused1 was removed and never used.
	FeatureLevelUnused1

	// FeatureLevelStoresChunkDataShaHashes is when chunk SHA1 hashes were added for faster comparisons.
	FeatureLevelStoresChunkDataShaHashes

	// FeatureLevelStoresPrerequisiteIDs is the version where prerequisite IDs were included.
	FeatureLevelStoresPrerequisiteIDs

	// FeatureLevelStoredAsBinaryData is when the first minimal binary format was introduced and UObject classes
	// were no longer saved in binary mode.
	FeatureLevelStoredAsBinaryData

	// FeatureLevelVariableSizeChunksWithoutWindowSizeChunkInfo is a temporary version where manifests referenced
	// variable-size chunks but didn’t serialize window size info. Elements move to ChunksV4.
	FeatureLevelVariableSizeChunksWithoutWindowSizeChunkInfo

	// FeatureLevelVariableSizeChunks is when variable-size chunks were fully supported, including serialization.
	FeatureLevelVariableSizeChunks

	// FeatureLevelUsesRuntimeGeneratedBuildID is when the manifest’s build ID started being generated from its metadata.
	FeatureLevelUsesRuntimeGeneratedBuildID

	// FeatureLevelUsesBuildTimeGeneratedBuildID is when a unique build-time generated ID began being stored in the
	// manifest.
	FeatureLevelUsesBuildTimeGeneratedBuildID

	// FeatureLevelLatestPlusOne is always one greater than the latest defined version.
	FeatureLevelLatestPlusOne

	// FeatureLevelLatest is an alias pointing to the actual latest feature level.
	FeatureLevelLatest = FeatureLevelLatestPlusOne - 1

	// FeatureLevelLatestNoChunks is the latest feature level supported by no-chunk (file data only) manifests.
	FeatureLevelLatestNoChunks = FeatureLevelStoresChunkFileSizes

	// FeatureLevelLatestJSON is the latest feature level supported by JSON-serialized manifests.
	FeatureLevelLatestJSON = FeatureLevelStoresPrerequisiteIDs

	// FeatureLevelFirstOptimisedDelta is the first version that supports optimized delta manifest generation.
	FeatureLevelFirstOptimisedDelta = FeatureLevelUsesRuntimeGeneratedBuildID

	// FeatureLevelStoresUniqueBuildID is an alias for FeatureLevelUsesRuntimeGeneratedBuildID.
	FeatureLevelStoresUniqueBuildID = FeatureLevelUsesRuntimeGeneratedBuildID

	// FeatureLevelBrokenJSONVersion is the JSON-manifest bug version (255), treated as FeatureLevelStoresChunkFileSizes.
	FeatureLevelBrokenJSONVersion = 255

	// FeatureLevelInvalid is the UObject default value, ensuring serialization always happens.
	FeatureLevelInvalid = -1
)

// ChunkSubDir returns the Chunk subdirectory name for the given FeatureLevel.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ManifestData.cpp#L78-L84
func (fl FeatureLevel) ChunkSubDir() string {
	switch {
	case fl > FeatureLevelStoredAsBinaryData:
		return "ChunksV4"
	case fl > FeatureLevelStoresDataGroupNumbers:
		return "ChunksV3"
	case fl > FeatureLevelStartStoringVersion:
		return "ChunksV2"
	default:
		return "Chunks"
	}
}

// StoredAs represents how a Chunk is stored.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/684b4c133ed87e8050d1fdaa287242f0fe2c1153/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ChunkData.h#L20-L29
type StoredAs uint8

const (
	// StoredAsPlainText indicates the chunk data is stored uncompressed and unencrypted.
	StoredAsPlainText StoredAs = 0x00

	// StoredAsCompressed indicates the chunk data is compressed.
	StoredAsCompressed StoredAs = 0x01

	// StoredAsEncrypted indicates the chunk data is encrypted. If the compressed flag is also set, the data is decrypted before decompression.
	StoredAsEncrypted StoredAs = 0x02
)

// ChunkHeaderVersion describes ChunkHeader.Version.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ChunkData.cpp#L72-L82
type ChunkHeaderVersion uint8

const (
	// ChunkHeaderVersionInvalid is an invalid chunk version.
	ChunkHeaderVersionInvalid ChunkHeaderVersion = iota

	// ChunkHeaderVersionOriginal is the original chunk version.
	ChunkHeaderVersionOriginal

	// ChunkHeaderVersionStoresShaAndHashType is the chunk version that stores SHA and hash type.
	ChunkHeaderVersionStoresShaAndHashType

	// ChunkHeaderVersionStoresDataSizeUncompressed is the chunk version that stores uncompressed data size.
	ChunkHeaderVersionStoresDataSizeUncompressed

	// ChunkHeaderVersionLatestPlusOne is always one greater than the latest defined version.
	ChunkHeaderVersionLatestPlusOne

	// ChunkHeaderVersionLatest is an alias pointing to the actual latest chunk version.
	ChunkHeaderVersionLatest = ChunkHeaderVersionLatestPlusOne - 1
)

// ManifestMetaVersion describes ManifestMeta.DataVersion.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ManifestData.cpp#L42-L50
type ManifestMetaVersion uint8

const (
	// ManifestMetaVersionOriginal is the original manifest meta version.
	ManifestMetaVersionOriginal ManifestMetaVersion = iota

	// ManifestMetaVersionSerialisesBuildID is the manifest meta version where the build ID started being serialized in the manifest metadata.
	ManifestMetaVersionSerialisesBuildID

	// ManifestMetaVersionLatestPlusOne is always one greater than the latest defined version.
	ManifestMetaVersionLatestPlusOne

	// ManifestMetaVersionLatest is an alias pointing to the actual latest manifest meta version.
	ManifestMetaVersionLatest = ManifestMetaVersionLatestPlusOne - 1
)

// ManifestHeaderVersionSizes contains the header size in bytes for each feature level.
// Versions up to FeatureLevelStoresPrerequisiteIDs use 37 bytes; versions from FeatureLevelStoredAsBinaryData onward use 41 bytes.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ManifestData.cpp#L28-L36
var ManifestHeaderVersionSizes = [FeatureLevelLatestPlusOne]int32{
	// FeatureLevelOriginal through FeatureLevelStoresPrerequisiteIDs: 37 bytes
	// (32b Magic, 32b HeaderSize, 32b DataSizeUncompressed, 32b DataSizeCompressed, 160b SHA1, 8b StoredAs)
	37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37,
	// FeatureLevelStoredAsBinaryData through FeatureLevelUsesBuildTimeGeneratedBuildID: 41 bytes
	// (296b Original, 32b Version)
	41, 41, 41, 41, 41,
}

// FileManifestListVersion describes FileManifestList.DataVersion.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ManifestData.cpp#L67-L74
type FileManifestListVersion uint8

const (
	// FileManifestListVersionOriginal is the original file manifest list version.
	FileManifestListVersionOriginal FileManifestListVersion = iota

	// FileManifestListVersionLatestPlusOne is always one greater than the latest defined version.
	FileManifestListVersionLatestPlusOne

	// FileManifestListVersionLatest is an alias pointing to the actual latest file manifest list version.
	FileManifestListVersionLatest = FileManifestListVersionLatestPlusOne - 1
)

// FileMetaFlag represents the EFileMetaFlags bitfield for file metadata flags.
//
// Source: https://github.com/EpicGames/UnrealEngine/blob/df42801f6a266711e1641d0058da4b0a0df711eb/Engine/Source/Runtime/Online/BuildPatchServices/Private/Data/ManifestData.h#L127-L136
type FileMetaFlag uint8

const (
	// FileMetaFlagNone indicates no special storage flags.
	FileMetaFlagNone FileMetaFlag = 0

	// FileMetaFlagReadOnly indicates the file is read-only.
	FileMetaFlagReadOnly FileMetaFlag = 1

	// FileMetaFlagCompressed indicates the file is compressed.
	FileMetaFlagCompressed FileMetaFlag = 1 << 1

	// FileMetaFlagUnixExecutable indicates the file is a Unix executable.
	FileMetaFlagUnixExecutable FileMetaFlag = 1 << 2
)
