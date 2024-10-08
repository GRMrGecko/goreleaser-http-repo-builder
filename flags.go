package main

import (
	"fmt"
	"github.com/alecthomas/kong"
)

type VersionFlag bool

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(appName + ": " + appVersion)
	app.Exit(0)
	return nil
}

// Flags supplied to cli.
type Flags struct {
	Version    VersionFlag   `name:"version" help:"Print version information and quit"`
	Repo       string        `help:"The path to a repo" required:"" type:"existingdir"`
	AddRelease AddReleaseCmd `cmd:"" help:"Add an release to the repo"`
	Prune      PruneCmd      `cmd:"" help:"Prune releases from repo."`
}

// Parse the supplied flags.
func (a *App) ParseFlags() *kong.Context {
	app.flags = &Flags{}

	ctx := kong.Parse(app.flags,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	return ctx
}
