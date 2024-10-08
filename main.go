package main

import (
	"time"
)

const (
	appName        = "goreleaser-http-repo-builder"
	appDescription = "Builds a repo for use with go-selfupdate"
	appVersion     = "0.1.0"
)

// App is the global application structure for communicating between servers and storing information.
type App struct {
	flags *Flags
	now   time.Time
}

var app *App

func main() {
	app = new(App)
	app.now = time.Now()
	ctx := app.ParseFlags()

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
