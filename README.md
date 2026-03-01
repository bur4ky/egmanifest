# Epic Games Manifest Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/bur4ky/egmanifest.svg)](https://pkg.go.dev/github.com/bur4ky/egmanifest)
[![CI](https://github.com/bur4ky/egmanifest/actions/workflows/ci.yml/badge.svg)](https://github.com/bur4ky/egmanifest/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/bur4ky/egmanifest)](LICENSE)

A Go library for parsing Epic Games manifest files  
This repository is a fork of [github.com/meszmate/manifest](https://github.com/meszmate/manifest) with many breaking changes, performance improvements and some features.

## Installation

```sh
go get github.com/bur4ky/egmanifest
```

## Basic Example

```go
package main

import (
	"log"
	"github.com/bur4ky/egmanifest"
)

const baseURL = "http://epicgames-download1.akamaized.net/Builds/Fortnite/CloudDir/"

func main() {
	manifestBytes := getManifestFile() // your own function
	binary, err := egmanifest.Parse(manifestBytes)
	if err != nil {
		log.Fatalln("Failed to parse manifest:", err)
	}

	for _, fm := range binary.FileManifestList.Files {
		for _, cp := range fm.ChunkParts {
			url := cp.Chunk.URL(baseURL + binary.Meta.FeatureLevel.ChunkSubDir())
			log.Println("Chunk URL:", url)

			// download, parse (manifest.ParseChunk() is recommended) and verify the chunk here
		}
	}
}

```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
