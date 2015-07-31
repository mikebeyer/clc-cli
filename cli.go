package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/aa"
	"github.com/mikebeyer/clc-cli/alert"
	"github.com/mikebeyer/clc-cli/group"
	"github.com/mikebeyer/clc-cli/lb"
	"github.com/mikebeyer/clc-cli/server"
	"github.com/mikebeyer/clc-cli/status"
	"github.com/mikebeyer/clc-sdk"
	"github.com/mikebeyer/clc-sdk/api"
)

func main() {

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
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "default-config", Usage: "create template configuration file"},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("default-config") {

			conf := api.NewConfig("USERNAME", "PASSWORD", "DEFAULT-ALIAS")
			b, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				fmt.Printf("unable to generate config template.")
			}

			err = ioutil.WriteFile("./clc.json", b, 0666)
			if err != nil {
				fmt.Printf("unable to generate config template.")
			}
			fmt.Printf("config template written to clc.json")
			return
		} else if !c.Args().Present() {
			cli.ShowAppHelp(c)
		}
	}
	var config api.Config
	config, err := api.FileConfig("./clc.json")
	if err != nil {
		config = api.EnvConfig()
	}

	client := clc.New(config)
	app.Commands = []cli.Command{
		server.Commands(client),
		status.Commands(client),
		aa.Commands(client),
		alert.Commands(client),
		lb.Commands(client),
		group.Commands(client),
	}
	app.Run(os.Args)
}
