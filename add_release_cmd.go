package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AddReleaseCmd struct {
	Release        string    `help:"Path to goreleaser dist folder." required:"" type:"existingdir"`
	Notes          string    `help:"Notes about this release."`
	Draft          bool      `help:"Is this release a draft?"`
	Prerelease     bool      `help:"Is this a prelease?"`
	IncludeBinary  bool      `help:"Include binary artifacts."`
	Force          bool      `help:"Force add, removing existing if needed."`
	PublishedAt    time.Time `help:"Specify exact time for release."`
	PublishedAtNow bool      `help:"Use the current time for published at instead of the metadata date."`
}

// Adds a release to a repo.
func (a *AddReleaseCmd) Run() error {
	// Read existing manifest for repo.
	manifestFile := filepath.Join(app.flags.Repo, "manifest.yaml")
	manifest, err := readManifestFile(manifestFile)
	if os.IsNotExist(err) {
		err = os.MkdirAll(app.flags.Repo, 0755)
	}
	if err != nil {
		return err
	}

	// Update old releases to include the ID field.
	for i, release := range manifest.Releases {
		manifest.Releases[i].ID = release.ReleaseID
	}

	// Read metadata from goreleaser.
	metadata, err := readMetadataFile(filepath.Join(a.Release, "metadata.json"))
	if err != nil {
		return err
	}
	versionPath := filepath.Join(app.flags.Repo, metadata.Version)

	// Read the artifcats to ensure we have a valid release.
	artifacts, err := readArtifactFile(filepath.Join(a.Release, "artifacts.json"))
	if err != nil {
		return err
	}
	if len(artifacts) == 0 {
		return errors.New("no artifacts in release")
	}

	// Validate the base dir for artifacts. It could be one dir up, or 2 dirs up.
	artifcatBase := a.Release
	artifactLayers := 0
	if _, serr := os.Stat(filepath.Join(artifcatBase, artifacts[0].Path)); serr != nil {
		artifcatBase = filepath.Dir(artifcatBase)
		artifactLayers = 1
		if _, serr := os.Stat(filepath.Join(artifcatBase, artifacts[0].Path)); serr != nil {
			artifcatBase = filepath.Dir(artifcatBase)
			artifactLayers = 2
			if _, serr := os.Stat(filepath.Join(artifcatBase, artifacts[0].Path)); serr != nil {
				return errors.New("unable to determine artificate base path")
			}
		}
	}

	// Check if the version already exists.
	existingIndex := -1
	for i, release := range manifest.Releases {
		if release.TagName == metadata.Version {
			existingIndex = i
			break
		}
	}

	// If the version already exists, ask about replacing.
	if existingIndex != -1 {
		if !a.Force {
			ans := askForConfirmation("This release already exists, should we replace?")

			// If we don't want to replace, we should stop here.
			if !ans {
				return errors.New("version already exists")
			}
		}

		// We need to replace the release, so remove it.
		manifest.Releases = append(manifest.Releases[:existingIndex], manifest.Releases[existingIndex+1:]...)

		// Remove the version directory.
		os.RemoveAll(versionPath)
	}

	// Make the release.
	manifest.LastReleaseID++
	release := &HttpRelease{
		ID:           manifest.LastReleaseID,
		ReleaseID:    manifest.LastReleaseID,
		Name:         metadata.Name,
		TagName:      metadata.Version,
		URL:          metadata.Version,
		Draft:        a.Draft,
		Prerelease:   a.Prerelease,
		PublishedAt:  metadata.Date,
		ReleaseNotes: a.Notes,
	}

	// If the publish date provided is valid, override.
	if !a.PublishedAt.IsZero() {
		release.PublishedAt = a.PublishedAt
	}

	// If published at is requested to be now, override.
	if a.PublishedAtNow {
		release.PublishedAt = app.now
	}

	// Make the directory for the release.
	err = os.Mkdir(versionPath, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Error making release directory: %s", err)
	}

	// Add artifacts.
	for _, artifact := range artifacts {
		// Skip binaries if not included.
		if artifact.Type == "Binary" && !a.IncludeBinary {
			continue
		}

		// Get the file path and confirm it exists and get its stat for file size.
		path := filepath.Join(artifcatBase, artifact.Path)
		stat, serr := os.Stat(path)
		if serr != nil {
			log.Println("Ignoring artifact", artifact.Name, "as its file does not exist.")
			continue
		}

		// Determine relative path.
		s := strings.Split(artifact.Path, "/")
		relativePath := filepath.Join(s[artifactLayers:]...)

		// Determine if artifact is in its own sub dir, make sure it exists.
		dir := filepath.Dir(relativePath)
		if dir != "." {
			os.MkdirAll(filepath.Join(versionPath, dir), 0755)
		}

		// Copy artifact to repo.
		err = copyFile(path, filepath.Join(versionPath, relativePath))
		if err != nil {
			log.Printf("Failed to copy artifact, skipping it: %s", err)
			continue
		}

		// Make asset.
		manifest.LastAssetID++
		asset := &HttpAsset{
			ID:   manifest.LastAssetID,
			Name: artifact.Name,
			Size: int(stat.Size()),
			URL:  filepath.Join(metadata.Version, relativePath),
		}

		// Add to the release.
		release.Assets = append(release.Assets, asset)
	}

	// Add release to manifest.
	manifest.Releases = append(manifest.Releases, release)

	// Write the manifest.
	err = writeManifestFile(manifestFile, manifest)
	if err != nil {
		return err
	}

	// If not a draft or prerelease, link latest to this release.
	if !a.Draft && !a.Prerelease {
		latestPath := filepath.Join(app.flags.Repo, "latest")
		os.Remove(latestPath)
		os.Symlink(metadata.Version, latestPath)
	}

	log.Println("Added release", metadata.Version, "for", metadata.Name, "to the repo", app.flags.Repo)

	return nil
}
