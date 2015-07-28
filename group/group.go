package group

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/util"
	"github.com/mikebeyer/clc-sdk"
	"github.com/mikebeyer/clc-sdk/group"
)

func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "group",
		Aliases: []string{"g"},
		Usage:   "group api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
			delete(client),
		},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get group details",
		Before: func(c *cli.Context) error {
			return util.CheckArgs(c)
		},
		Action: func(c *cli.Context) {
			resp, err := client.Group.Get(c.Args().First())
			if err != nil {
				fmt.Printf("failed to get %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func create(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create group",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "group name [required]"},
			cli.StringFlag{Name: "parent, p", Usage: "parent group id [required]"},
			cli.StringFlag{Name: "description, d", Usage: "group description"},
		},
		Before: func(c *cli.Context) error {
			return util.CheckStringFlag(c, "name", "parent")
		},
		Action: func(c *cli.Context) {
			g := group.Group{
				Name:          c.String("name"),
				ParentGroupID: c.String("parent"),
				Description:   c.String("description"),
			}
			resp, err := client.Group.Create(g)
			if err != nil {
				fmt.Printf("failed to create %s", c.String("name"))
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func delete(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete group details",
		Before: func(c *cli.Context) error {
			return util.CheckArgs(c)
		},
		Action: func(c *cli.Context) {
			resp, err := client.Group.Delete(c.Args().First())
			if err != nil {
				fmt.Printf("failed to delete %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}
