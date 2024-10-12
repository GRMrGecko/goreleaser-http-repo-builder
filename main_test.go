package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test the add release option.
func TestAppFunctionality(t *testing.T) {
	// Make temp directory to build repo.
	dname, err := os.MkdirTemp("", "goreleaser-http-repo-builder")
	if err != nil {
		t.Errorf("error making tempdir: %s", err)
	}

	// Get the tests dir with test files.
	testsDir, err := filepath.Abs("tests")
	if err != nil {
		t.Errorf("error finding tests dir: %s", err)
	}

	// Now date for app defines.
	now, _ := time.Parse(time.DateOnly, "2024-10-08")

	// Test adding a release of v0.1.
	os.Args = []string{"test", "--repo", dname, "add-release", "--notes", "This is a test.", "--release", filepath.Join(testsDir, "v0.1")}
	app = new(App)
	app.now = now
	ctx := app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Test adding a release of v0.1.1.
	os.Args = []string{"test", "--repo", dname, "add-release", "--draft", "--release", filepath.Join(testsDir, "v0.1.1")}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Test adding a release of v0.1.2.
	os.Args = []string{"test", "--repo", dname, "add-release", "--prerelease", "--published-at", "2024-10-05T22:15:21.731224367-05:00", "--release", filepath.Join(testsDir, "v0.1.2")}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Confirm the latest release is v0.1.0.
	latestPath, err := os.Readlink(filepath.Join(dname, "latest"))
	if err != nil {
		t.Errorf("error reading link to latest: %s", err)
	}
	if latestPath != "v0.1.0" {
		t.Error("the latest link isn't correctly linked")
	}

	// Hash the manifest file without the published_at dates.
	hfun := md5.New()
	d, err := os.ReadFile(filepath.Join(dname, "manifest.yaml"))
	if err != nil {
		t.Errorf("error reading manifest file: %s", err)
	}

	// Hash the result and confirm.
	hfun.Write(d)
	sum := hfun.Sum(nil)
	hash := hex.EncodeToString(sum)
	if hash != "19a3a502913252635b3e0ea838846197" {
		t.Errorf("hash isn't valid for manifest file: %s", hash)
	}

	// Binaries are not included.
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.0/example_linux_amd64/example")); !os.IsNotExist(serr) {
		t.Error("v0.1.0 binary exists, when it shouldn't exist.")
	}

	// Confirm the asset was copied correctly.
	d, err = os.ReadFile(filepath.Join(dname, "v0.1.0/example_linux_amd64.tar.gz"))
	if err != nil {
		t.Errorf("error reading test file: %s", err)
	}
	hfun.Reset()
	hfun.Write(d)
	sum = hfun.Sum(nil)
	hash = hex.EncodeToString(sum)
	if hash != "9cffcbe826ae684db1c8a08ff9216f34" {
		t.Errorf("hash isn't valid for test file: %s", hash)
	}

	// Confirm pruning of max releases works.
	os.Args = []string{"test", "--repo", dname, "prune", "--max-age=216h"}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Confirm pruned state.
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.0/example_linux_amd64.tar.gz")); !os.IsNotExist(serr) {
		t.Error("v0.1.0 exists, when it shouldn't exist.")
	}
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.1/example_linux_amd64.tar.gz")); os.IsNotExist(serr) {
		t.Error("v0.1.1 does not exists, when it should.")
	}
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.2/example_linux_amd64.tar.gz")); os.IsNotExist(serr) {
		t.Error("v0.1.2 does not exists, when it should.")
	}

	// Delete all files, and reset.
	os.RemoveAll(dname)
	os.Mkdir(dname, 0755)

	// Test adding a release of v0.1.
	os.Args = []string{"test", "--repo", dname, "add-release", "--include-binary", "--release", filepath.Join(testsDir, "v0.1")}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Test adding a release of v0.1.1.
	os.Args = []string{"test", "--repo", dname, "add-release", "--release", filepath.Join(testsDir, "v0.1.1")}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Test adding a release of v0.1.2.
	os.Args = []string{"test", "--repo", dname, "add-release", "--published-at-now", "--release", filepath.Join(testsDir, "v0.1.2")}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Confirm the latest release is v0.1.2.
	latestPath, err = os.Readlink(filepath.Join(dname, "latest"))
	if err != nil {
		t.Errorf("error reading link to latest: %s", err)
	}
	if latestPath != "v0.1.2" {
		t.Error("the latest link isn't correctly linked")
	}

	// Hash the manifest file without the published_at dates.
	hfun.Reset()
	d, err = os.ReadFile(filepath.Join(dname, "manifest.yaml"))
	if err != nil {
		t.Errorf("error reading manifest file: %s", err)
	}

	// Hash the result and confirm.
	hfun.Write(d)
	sum = hfun.Sum(nil)
	hash = hex.EncodeToString(sum)
	if hash != "999c4156c2b5ff25f3491b86c8255cb5" {
		t.Errorf("hash isn't valid for manifest file: %s", hash)
	}

	// Binaries are not included.
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.0/example_linux_amd64/example")); os.IsNotExist(serr) {
		t.Error("v0.1.0 binary does not exists, when it shouldn.")
	}

	// Confirm pruning of max releases works.
	os.Args = []string{"test", "--repo", dname, "prune", "--max-releases=1"}
	app = new(App)
	app.now = now
	ctx = app.ParseFlags()

	// Run the command.
	err = ctx.Run()
	if err != nil {
		t.Errorf("error running the app: %s", err)
	}

	// Confirm pruned state.
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.0/example_linux_amd64.tar.gz")); !os.IsNotExist(serr) {
		t.Error("v0.1.0 exists, when it shouldn't exist.")
	}
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.1/example_linux_amd64.tar.gz")); !os.IsNotExist(serr) {
		t.Error("v0.1.1 exists, when it shouldn't exist.")
	}
	if _, serr := os.Stat(filepath.Join(dname, "v0.1.2/example_linux_amd64.tar.gz")); os.IsNotExist(serr) {
		t.Error("v0.1.2 does not exists, when it should.")
	}

	// Cleanup.
	os.RemoveAll(dname)
}
