package util

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

func CheckStringFlag(c *cli.Context, flags ...string) error {
	for _, v := range flags {
		if c.String(v) == "" {
			return displayAndErr(fmt.Sprintf("--%s is required [--help for usage]", v))
		}
	}
	return nil
}

func CheckStringSliceFlag(c *cli.Context, flag string) error {
	if len(c.StringSlice(flag)) == 0 {
		return displayAndErr(fmt.Sprintf("at least once --%s is required [--help for usage]", flag))
	}
	return nil
}

func CheckForAnyOfStringSliceFlag(c *cli.Context, flags ...string) error {
	accum := 0
	for _, v := range flags {
		accum += len(c.StringSlice(v))
	}
	if accum == 0 {
		return displayAndErr(fmt.Sprintf("at least one of --%s is required [--help for usage]", strings.Join(flags, ", --")))
	}
	return nil
}

func CheckArgs(c *cli.Context) error {
	if !c.Args().Present() {
		return displayAndErr(fmt.Sprintf("missing arguments [--help for usage]"))
	}
	return nil
}

func displayAndErr(msg string) error {
	fmt.Println(msg)
	return fmt.Errorf(msg)
}
