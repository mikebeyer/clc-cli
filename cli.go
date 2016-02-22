package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/api"
	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/aa"
	"github.com/mikebeyer/clc-cli/alert"
	"github.com/mikebeyer/clc-cli/group"
	"github.com/mikebeyer/clc-cli/lb"
	"github.com/mikebeyer/clc-cli/server"
	"github.com/mikebeyer/clc-cli/status"
)

func main() {
	app := cli.NewApp()
	app.Name = "clc"
	app.Usage = "clc v2 api cli"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mike Beyer",
			Email: "michael.beyer@ctl.io",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "gen-config", Usage: "create template configuration file"},
		cli.StringFlag{Name: "config", Usage: "specify config file"},
	}
	var config api.Config
	var configErr error
	app.Action = func(c *cli.Context) {
		if c.Bool("gen-config") {
			conf, err := api.NewConfig("USERNAME", "PASSWORD", "DEFAULT-ALIAS", "")
			if err != nil {
				fmt.Printf("unable to generate config template.\n")
				return
			}
			b, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				fmt.Printf("unable to generate config template.\n")
				return
			}

			err = ioutil.WriteFile("./config.json", b, 0666)
			if err != nil {
				fmt.Printf("unable to generate config template.\n")
				return
			}
			fmt.Printf("config template written to config.json\n")
			return
		} else if c.Bool("help") {
			cli.ShowAppHelp(c)
		} else if !c.Args().Present() {
			if config.Valid() || c.GlobalString("config") != "" {
				cli.ShowAppHelp(c)
			} else {
				if !config.Valid() && c.GlobalString("config") == "" {
					config, configErr = api.EnvConfig()
					if configErr != nil {
						config, configErr = api.FileConfig("./config.json")
						if configErr != nil {
							fmt.Printf("WARNING: failed to find necessary environment variables or default config location (./config.json)\n")
							fmt.Printf("WARNING: --configure to generate configuration file\n")
							fmt.Printf("WARNING: --help for all options\n")
							return
						}
					}
					cli.ShowAppHelp(c)
				}
			}
		}
	}
	config, _ = api.EnvConfig()

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
