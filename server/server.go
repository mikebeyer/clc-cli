package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/server"
	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/util"
)

// Commands exports the cli commands for the server package
func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "server api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
			delete(client),
			publicIP(client),
			archive(client),
			restore(client),
			sshServer(client),
		},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get server details",
		Before: func(c *cli.Context) error {
			return util.CheckArgs(c)
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			server, err := client.Server.Get(c.Args().First())
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func createFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
		cli.StringFlag{Name: "cpu, c", Usage: "server cpus (1 - 16) [required]"},
		cli.StringFlag{Name: "memory, m", Usage: "server memory in gbs (1 - 128) [required]"},
		cli.StringFlag{Name: "group, g", Usage: "parent group id [required]"},
		cli.StringFlag{Name: "source, s", Usage: "source server id (template or existing server) [required or helper flag]"},
		cli.BoolFlag{Name: "standard", Usage: "standard server"},
		cli.BoolFlag{Name: "hyperscale", Usage: "hyperscale server [overides storage setting]"},
		cli.BoolFlag{Name: "premium", Usage: "premium storage"},
		cli.StringFlag{Name: "password, p", Usage: "server password"},
		cli.StringFlag{Name: "description, d", Usage: "server description"},
		cli.StringFlag{Name: "ip", Usage: "id address"},
		cli.BoolFlag{Name: "managed", Usage: "make server managed"},
		cli.StringFlag{Name: "primaryDNS", Usage: "primary dns"},
		cli.StringFlag{Name: "secondaryDNS", Usage: "secondary dns"},
		cli.StringFlag{Name: "network", Usage: "network id"},
		cli.StringFlag{Name: "storage", Usage: "standard or premium"},
	}
}

func deriveServerType(c *cli.Context) string {
	if c.Bool("standard") {
		return "standard"
	}
	return "hyperscale"
}

func deriveStorageType(c *cli.Context) string {
	if c.Bool("premium") && c.Bool("standard") {
		return "premium"
	} else if c.Bool("standard") {
		return "standard"
	}
	return "hyperscale"
}

func create(client *clc.Client) cli.Command {
	return cli.Command{
		Name:        "create",
		Aliases:     []string{"c"},
		Usage:       "create server",
		Description: "example: clc server create -n 'my cool server' -c 1 -m 1 -g [group id] -t standard --ubuntu-14",
		Flags:       append(createFlags(), templateFlags()...),
		Before: func(c *cli.Context) error {
			err := util.CheckStringFlag(c, "name", "cpu", "memory", "group")
			if err == nil {
				err = util.CheckForEitherBooleanFlag(c, "standard", "hyperscale")
			}
			return err
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			source, err := findTemplateInContext(c)
			if err != nil {
				return
			}
			server := server.Server{
				Name:           c.String("name"),
				CPU:            c.Int("cpu"),
				MemoryGB:       c.Int("memory"),
				GroupID:        c.String("group"),
				SourceServerID: source,
				Type:           deriveServerType(c),
			}
			server.Password = c.String("password")
			server.Description = c.String("description")
			server.IPaddress = c.String("ip")
			server.IsManagedOS = c.Bool("managed")
			server.PrimaryDNS = c.String("primaryDNS")
			server.SecondaryDNS = c.String("secondaryDNS")
			server.NetworkID = c.String("network")
			server.Storagetype = deriveStorageType(c)

			resp, err := client.Server.Create(server)
			if err != nil {
				fmt.Printf("failed to create %s", server.Name)
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
		Usage:   "delete server",
		Before: func(c *cli.Context) error {
			return util.CheckArgs(c)
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			server, err := client.Server.Delete(c.Args().First())
			if err != nil {
				fmt.Printf("failed to delete %s", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				fmt.Printf("%s", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func archive(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "archive",
		Aliases: []string{"a"},
		Usage:   "archive server",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "name, n", Usage: "name of servers to archive"},
		},
		Before: func(c *cli.Context) error {
			return util.CheckStringSliceFlag(c, "name")
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp, err := client.Server.Archive(c.StringSlice("name")...)
			if err != nil {
				fmt.Printf("failed to archive %s", strings.Join(c.StringSlice("name"), ", "))
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

func restore(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "restore",
		Aliases: []string{"r"},
		Usage:   "restore server",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "name of server to restore [required]"},
			cli.StringFlag{Name: "group, g", Usage: "group for server to restore to [required]"},
		},
		Before: func(c *cli.Context) error {
			return util.CheckStringFlag(c, "name", "group")
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp, err := client.Server.Restore(c.String("name"), c.String("group"))
			if err != nil {
				fmt.Printf("failed to restore %s", c.String("name"))
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

func publicIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "public-ip",
		Aliases: []string{"ip"},
		Usage:   "manage public ips",
		Subcommands: []cli.Command{
			getIP(client),
			createIP(client),
			deleteIP(client),
		},
	}
}

func getIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get public ip",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringFlag{Name: "ip", Usage: "ip [required]"},
		},
		Before: func(c *cli.Context) error {
			return util.CheckStringFlag(c, "name", "ip")
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp, err := client.Server.GetPublicIP(c.String("name"), c.String("ip"))
			if err != nil {
				fmt.Printf("err %s\n", err)
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

func createIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add public ip to server",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringSliceFlag{Name: "tcp", Usage: "provide a port [8080] or a port range [8080:8082]"},
			cli.StringSliceFlag{Name: "udp", Usage: "provide a port [8080] or a port range [8080:8082]"},
			cli.StringSliceFlag{Name: "restriction, r", Usage: "provide an ip subnet to restrict to access the public ip [ex. 10.0.0.1/24 (must be cidr notation)]"},
		},
		Before: func(c *cli.Context) error {
			err := util.CheckStringFlag(c, "name")
			if err == nil {
				err = util.CheckForAnyOfStringSliceFlag(c, "tcp", "udp")
			}
			return err
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ports := make([]server.Port, 0)
			tcps, err := parsePort("tcp", c.StringSlice("tcp"))
			if err != nil {
				fmt.Println(err.Error())
			}
			ports = append(ports, tcps...)
			udps, err := parsePort("udp", c.StringSlice("udp"))
			if err != nil {
				fmt.Println(err.Error())
			}
			ports = append(ports, udps...)
			restrictions := make([]server.SourceRestriction, 0)
			for _, v := range c.StringSlice("restriction") {
				restrictions = append(restrictions, server.SourceRestriction{CIDR: v})
			}

			ip := server.PublicIP{Ports: ports}
			resp, err := client.Server.AddPublicIP(c.String("name"), ip)
			if err != nil {
				fmt.Printf("err %s\n", err)
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

func parsePort(protocol string, list []string) ([]server.Port, error) {
	ports := make([]server.Port, 0)
	for _, v := range list {
		r := strings.Split(v, ":")
		port, err := strconv.Atoi(r[0])
		if err != nil {
			return ports, fmt.Errorf("invalid port provided %s", r[0])
		}
		if len(r) > 1 {
			to, err := strconv.Atoi(r[1])
			if err != nil {
				return ports, fmt.Errorf("invalid port provided %s", r[0])
			}
			ports = append(ports, server.Port{Protocol: protocol, Port: port, PortTo: to})
		} else {
			ports = append(ports, server.Port{Protocol: protocol, Port: port})
		}
	}
	return ports, nil
}

func deleteIP(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete public ip",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "server name [required]"},
			cli.StringFlag{Name: "ip", Usage: "ip [required]"},
		},
		Before: func(c *cli.Context) error {
			return util.CheckStringFlag(c, "name", "ip")
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp, err := client.Server.DeletePublicIP(c.String("name"), c.String("ip"))
			if err != nil {
				fmt.Printf("err %s\n", err)
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

func sshServer(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "ssh",
		Aliases: []string{"s"},
		Usage:   "ssh to server",
		Before: func(c *cli.Context) error {
			return util.CheckArgs(c)
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			server, err := client.Server.Get(c.Args().First())
			if err != nil {
				fmt.Printf("failed to shh to %s\n", c.Args().First())
				return
			}
			creds, err := client.Server.GetCredentials(c.Args().First())
			if err != nil {
				fmt.Printf("failed to shh to %s\n", c.Args().First())
				return
			}

			internal, err := getInternal(server)
			if err != nil {
				fmt.Printf("failed to shh to %s\n", c.Args().First())
				return
			}

			if err := connect(internal, creds); err != nil {
				fmt.Printf("failed to shh to %s\n", c.Args().First())
				return
			}
		},
	}
}

func getInternal(server *server.Response) (string, error) {
	for _, ip := range server.Details.IPaddresses {
		if ip.Public == "" && ip.Internal != "" {
			return ip.Internal, nil
		}
	}
	return "", fmt.Errorf("unable to find ip")
}
