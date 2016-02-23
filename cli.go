package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

var client *clc.Client

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
		cli.BoolFlag{Name: "gen-config-from-env", Usage: "create configuration file from environment variables"},
	}
	app.Action = func(c *cli.Context) {
		if !c.Bool("help") {
			if c.Bool("gen-config") {
				err := genConfig()
				if err != nil {
					fmt.Println("failed to generate default config")
				}
				return
			}

			if c.Bool("gen-config-from-env") {
				err := genConfigFromEnv()
				if err != nil {
					fmt.Println("failed to generate default config")
				}
				return
			}

			err := loadClient()
			if err != nil && !c.Args().Present() {
				fmt.Println("WARNING: failed to find necessary environment variables or default config location (./config.json)")
				fmt.Println("         --gen-config to generate configuration file")
				fmt.Println("         --help for all options")
			} else {
				cli.ShowAppHelp(c)
			}
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

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func loadClient() error {
	var config api.Config

	config, err := api.FileConfig("./config.json")
	if err != nil {
		config, err = api.EnvConfig()
		if err != nil {

			return err
		}
	}

	client = clc.New(config)
	return err
}

func genConfigFromEnv() error {
	config, err := api.EnvConfig()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./config.json", b, 0666)
	if err != nil {
		return err
	}
	fmt.Printf("config written to config.json from current environment varibales\n")
	return nil
}

func genConfig() error {
	config, err := api.NewConfig("USERNAME", "PASSWORD", "DEFAULT-ALIAS", "")
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./config.json", b, 0666)
	if err != nil {
		return err
	}
	fmt.Printf("config template written to config.json\n")
	return nil
}
