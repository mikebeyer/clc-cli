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
		cli.BoolFlag{Name: "gen-config", Usage: "create template configuration file"},
	}
	var config api.Config
	var client *clc.Client
	var configErr error
	app.Action = func(c *cli.Context) {
		if c.Bool("gen-config") {
			conf := api.NewConfig("USERNAME", "PASSWORD", "DEFAULT-ALIAS")
			b, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				fmt.Printf("unable to generate config template.")
			}

			err = ioutil.WriteFile("./config.json", b, 0666)
			if err != nil {
				fmt.Printf("unable to generate config template.")
			}
			fmt.Printf("config template written to config.json")
			return
		} else if c.String("config") != "" {
			config, configErr = api.FileConfig(c.String("config"))
			if configErr != nil {
				fmt.Printf("unable to find/parse config: %s\n", c.String("config"))
			}
		} else if c.Bool("help") {
			cli.ShowAppHelp(c)
		} else if !c.Args().Present() {
			if !config.Valid() {
				config, configErr = api.EnvConfig()
				if configErr != nil {
					config, configErr = api.FileConfig("./config.json")
					if configErr != nil {
						fmt.Printf("failed to find necessary environment variables or default config location (./config.json)\n")
						return
					}
				}
			}
			client = clc.New(config)
		}
	}

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
