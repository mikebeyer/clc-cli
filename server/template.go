package server

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/util"
)

func templateFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{Name: "pxe", Usage: "pxe boot [experimental]"},
		cli.BoolFlag{Name: "centos-5", Usage: "centos 5"},
		cli.BoolFlag{Name: "centos-6", Usage: "centos 6"},
		cli.BoolFlag{Name: "debian-6", Usage: "debian 6"},
		cli.BoolFlag{Name: "debian-7", Usage: "debian 7"},
		cli.BoolFlag{Name: "rhel-5", Usage: "rhel 5"},
		cli.BoolFlag{Name: "rhel-6", Usage: "rhel 6"},
		cli.BoolFlag{Name: "rhel-7", Usage: "rhel 7"},
		cli.BoolFlag{Name: "ubuntu-12", Usage: "ubuntu 12"},
		cli.BoolFlag{Name: "ubuntu-14", Usage: "ubuntu 14"},
		cli.BoolFlag{Name: "win2008r2-dc", Usage: "windows 2008 r2 datacenter edition"},
		cli.BoolFlag{Name: "win2008r2-ent", Usage: "windows 2008 r2 enterprise edition"},
		cli.BoolFlag{Name: "win2008r2-std", Usage: "windows 2008 r2 stanadard edition"},
		cli.BoolFlag{Name: "win2012", Usage: "windows 2012 datacenter edition"},
		cli.BoolFlag{Name: "win2012r2", Usage: "windows 2012 r2 datacenter edition"},
	}
}

var (
	templates = map[string]string{
		"pxe":           "PXE-TEMPLATE",
		"centos-5":      "CENTOS-5-64-TEMPLATE",
		"centos-6":      "CENTOS-6-64-TEMPLATE",
		"debian-6":      "DEBIAN-6-64-TEMPLATE",
		"debian-7":      "DEBIAN-7-64-TEMPLATE",
		"rhel-5":        "RHEL-5-64-TEMPLATE",
		"rhel-6":        "RHEL-6-64-TEMPLATE",
		"rhel-7":        "RHEL-7-64-TEMPLATE",
		"ubuntu-12":     "UBUNTU-12-64-TEMPLATE",
		"ubuntu-14":     "UBUNTU-14-64-TEMPLATE",
		"win2008r2-dc":  "WIN2008R2DTC-64",
		"win2008r2-ent": "WIN2008R2ENT-64",
		"win2008r2-std": "WIN2008R2STD-64",
		"win2012":       "WIN2012DTC-64",
		"win2012r2":     "WIN2012R2DTC-64",
	}
)

func findTemplateInContext(c *cli.Context) (string, error) {
	template := ""
	for k, v := range templates {
		if template != "" && c.Bool(k) {
			return "", util.DisplayAndErr(fmt.Sprintf("multiple templates provided. [--help for usage]"))
		}
		if c.Bool(k) {
			template = v
		}
	}

	if template != "" && c.String("source") != "" {
		return "", util.DisplayAndErr(fmt.Sprintf("multiple templates provided. [--help for usage]"))
	}

	if template == "" && c.String("source") != "" {
		template = c.String("source")
	}

	if template == "" {
		return "", util.DisplayAndErr(fmt.Sprintf("a template must be provided. [--help for usage]"))
	}
	return template, nil
}
