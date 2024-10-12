package main

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// An individual asset.
type HttpAsset struct {
	ID   int64  `yaml:"id"`
	Name string `yaml:"name"`
	Size int    `yaml:"size"`
	URL  string `yaml:"url"`
}

// An individual release.
type HttpRelease struct {
	ID           int64        `yaml:"id"`
	ReleaseID    int64        `yaml:"release_id"`
	Name         string       `yaml:"name"`
	TagName      string       `yaml:"tag_name"`
	URL          string       `yaml:"url"`
	Draft        bool         `yaml:"draft"`
	Prerelease   bool         `yaml:"prerelease"`
	PublishedAt  time.Time    `yaml:"published_at"`
	ReleaseNotes string       `yaml:"release_notes"`
	Assets       []*HttpAsset `yaml:"assets"`
}

// The manifest file structure.
type HttpManifest struct {
	LastReleaseID int64          `yaml:"last_release_id"`
	LastAssetID   int64          `yaml:"last_asset_id"`
	Releases      []*HttpRelease `yaml:"releases"`
}

// Read and parse manifest file.
func readManifestFile(manifestFile string) (*HttpManifest, error) {
	// We always want a manifest incase repo just needs to start from scratch.
	manifest := new(HttpManifest)

	// Read file, if error return the error.
	yamlFile, err := os.Open(manifestFile)
	if err != nil {
		return manifest, err
	}

	// Attempt to decode the file.
	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(manifest)
	yamlFile.Close()

	// Return the manifest and if any error occurred.
	return manifest, err
}

// Write manifest file.
func writeManifestFile(manifestFile string, manifest *HttpManifest) error {
	// Open the file for write.
	yamlFile, err := os.OpenFile(manifestFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer yamlFile.Close()

	// Encode data.
	encoder := yaml.NewEncoder(yamlFile)
	err = encoder.Encode(manifest)
	return err
}
