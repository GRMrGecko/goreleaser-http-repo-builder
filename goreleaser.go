package main

import (
	"encoding/json"
	"os"
	"time"
)

// The metadata needed from goreleaser.
type Metadata struct {
	Name    string    `json:"project_name"`
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

// Read and parse metadata file
func readMetadataFile(metadataFile string) (*Metadata, error) {
	// Read file, if error return the error.
	jsonFile, err := os.Open(metadataFile)
	if err != nil {
		return nil, err
	}

	// Attempt to decode the file.
	metadata := new(Metadata)
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(metadata)
	jsonFile.Close()

	// Return the metadata and if any error occurred.
	return metadata, err
}

// Artifcat map.
type Artifact struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// Read and parse metadata file
func readArtifactFile(artifactFile string) ([]*Artifact, error) {
	// Read file, if error return the error.
	jsonFile, err := os.Open(artifactFile)
	if err != nil {
		return nil, err
	}

	// Attempt to decode the file.
	var artifacts []*Artifact
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&artifacts)
	jsonFile.Close()

	// Return the metadata and if any error occurred.
	return artifacts, err
}
