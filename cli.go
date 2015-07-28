package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/aa"
	"github.com/mikebeyer/clc-cli/alert"
	"github.com/mikebeyer/clc-cli/lb"
	"github.com/mikebeyer/clc-cli/server"
	"github.com/mikebeyer/clc-cli/status"
	"github.com/mikebeyer/clc-sdk"
	"github.com/mikebeyer/clc-sdk/api"
)

func main() {
	client := clc.New(api.EnvConfig())

	app := cli.NewApp()
	app.Name = "clc"
	app.Usage = "clc v2 api cli"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mike Beyer",
			Email: "michael.beyer@ctl.io",
		},
	}
	app.Commands = []cli.Command{
		server.Commands(client),
		status.Commands(client),
		aa.Commands(client),
		alert.Commands(client),
		lb.Commands(client),
	}
	app.Run(os.Args)
}
