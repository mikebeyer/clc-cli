package util

import (
	"fmt"
	"strings"

	"github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/api"
	"github.com/codegangsta/cli"
)

func CheckStringFlag(c *cli.Context, flags ...string) error {
	for _, v := range flags {
		if c.String(v) == "" {
			return DisplayAndErr(fmt.Sprintf("--%s is required [--help for usage]", v))
		}
	}
	return nil
}

func CheckStringSliceFlag(c *cli.Context, flag string) error {
	if len(c.StringSlice(flag)) == 0 {
		return DisplayAndErr(fmt.Sprintf("at least once --%s is required [--help for usage]", flag))
	}
	return nil
}

func CheckForAnyOfStringSliceFlag(c *cli.Context, flags ...string) error {
	accum := 0
	for _, v := range flags {
		accum += len(c.StringSlice(v))
	}
	if accum == 0 {
		return DisplayAndErr(fmt.Sprintf("at least one of --%s is required [--help for usage]", strings.Join(flags, ", --")))
	}
	return nil
}

func CheckForEitherBooleanFlag(c *cli.Context, right, left string) error {
	if (c.Bool(right) && c.Bool(left)) || (!c.Bool(right) && !c.Bool(left)) {
		return DisplayAndErr(fmt.Sprintf("only one of --%s or --%s can be provided", right, left))
	}
	return nil
}

func CheckArgs(c *cli.Context) error {
	if !c.Args().Present() {
		return DisplayAndErr(fmt.Sprintf("missing arguments [--help for usage]"))
	}
	return nil
}

func DisplayAndErr(msg string) error {
	fmt.Println(msg)
	return fmt.Errorf(msg)
}

func MaybeLoadConfig(c *cli.Context, client *clc.Client) (*clc.Client, error) {
	if c.GlobalString("config") != "" {
		config, err := api.FileConfig(c.GlobalString("config"))
		if err != nil {
			return nil, fmt.Errorf("failed to load config at %s", c.GlobalString("config"))
		}

		return clc.New(config), nil
	}

	return client, nil
}
