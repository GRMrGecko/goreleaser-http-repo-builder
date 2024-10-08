package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type PruneCmd struct {
	MaxAge      time.Duration `help:"Delete releases older than."`
	MaxReleases int           `help:"Maximum number of releases to keep."`
	DryRun      bool          `help:"Just log the result without actually pruning."`
}

// Extra help to explain you can't set 2 prune stratages.
func (a *PruneCmd) Help() string {
	return "You cannot use both max-age and max-releases, only set one."
}

// Verify the options provided to the command.
func (a *PruneCmd) AfterApply() error {
	// If both stratages are defined, we don't allow that.
	if a.MaxAge > time.Duration(0) && a.MaxReleases > 0 {
		return errors.New("must only provide one prune argument")
	}
	// If no stratages are defined, we don't allow that.
	if a.MaxAge <= time.Duration(0) && a.MaxReleases <= 0 {
		return errors.New("must provide one prune argument")
	}
	return nil
}

// Adds a release to a repo.
func (a *PruneCmd) Run() error {
	// Read existing manifest for repo.
	manifestFile := filepath.Join(app.flags.Repo, "manifest.yaml")
	manifest, err := readManifestFile(manifestFile)
	if err != nil {
		return err
	}

	// Keep reference of number of pruned releases.
	releasesPruned := 0
	n := len(manifest.Releases)

	// If max releases defined and is less than number of releases, look for items to prune.
	if a.MaxReleases > 0 && n > a.MaxReleases {
		// Loop starting at max releases.
		for i := a.MaxReleases; i < n; i++ {
			// Get the current release.
			// We want pull from the top of the stack downward to keep newer releases.
			version := manifest.Releases[n-(i+1)].TagName
			log.Println("Removing release:", version)

			// If this isn't a dry run, remove the version directory.
			if !a.DryRun {
				err = os.RemoveAll(filepath.Join(app.flags.Repo, version))
				if err != nil {
					return fmt.Errorf("untable to remove release files: %s", err)
				}
			}

			// Count the number pruned.
			releasesPruned++
		}

		// Remove releases from the slice.
		manifest.Releases = manifest.Releases[n-a.MaxReleases:]
	}

	// If we are pruning based on duration, do so.
	if a.MaxAge > time.Duration(0) {
		// Loop through the releases, and find old releases.
		for i := 0; i < n; i++ {
			// If n is 1, we removed too many entries and need to stop.
			if n == 1 {
				log.Println("The repo has only 1 release remaining, ending the prune here to keep 1 release.")
				break
			}

			// Get the current release, and confirm its age.
			release := manifest.Releases[i]
			if app.now.Sub(release.PublishedAt) >= a.MaxAge {
				// This release is too old, so we need to remove it.
				version := release.TagName
				log.Println("Removing release:", version)

				// If this isn't a dry run, remove the version files.
				if !a.DryRun {
					err = os.RemoveAll(filepath.Join(app.flags.Repo, version))
					if err != nil {
						return fmt.Errorf("untable to remove release files: %s", err)
					}
				}

				// Remove the release from the slice.
				manifest.Releases = append(manifest.Releases[:i], manifest.Releases[i+1:]...)

				// Back up one on the index and number of releases as one release was deleted.
				i--
				n--

				// Count number of releases pruned.
				releasesPruned++

			}
		}
	}

	// Write the manifest if this isn't a dry run.
	if !a.DryRun {
		err = writeManifestFile(manifestFile, manifest)
		if err != nil {
			return err
		}
	}

	// Provide details on what's been pruned.
	log.Println("Pruned", releasesPruned, "release from the repo.")

	return nil
}
