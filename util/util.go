package util

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func CheckStringFlag(c *cli.Context, flags ...string) error {
	for _, v := range flags {
		if c.String(v) == "" {
			return fmt.Errorf("--%s is required [--help for usage]", v)
		}
	}
	return nil
}
