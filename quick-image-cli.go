package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "quick-image-cli"
	app.Version = Version
	app.Usage = "terminal tool to upload quickly and easily image"
	app.Author = "kyokomi"
	app.Email = "kyoko1220adword@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
