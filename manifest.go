// Package egmanifest parses Epic Games manifest files in both binary and JSON formats.
package egmanifest

// Parse detects whether `b` is JSON or binary, parses it, and returns a BinaryManifest.
func Parse(b []byte) (*BinaryManifest, error) {
	if len(b) > 0 && b[0] == '{' {
		jm, err := ParseJSON(b)
		if err != nil {
			return nil, err
		}

		return jm.ToBinaryManifest()
	}

	return ParseBinary(b)
}
