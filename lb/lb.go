package lb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/lb"
	"github.com/codegangsta/cli"
	"github.com/mikebeyer/clc-cli/util"
)

func Commands(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "load-balancer",
		Aliases: []string{"lb"},
		Usage:   "load balancer api",
		Subcommands: []cli.Command{
			get(client),
			create(client),
			update(client),
			delete(client),
			getPool(client),
			createPool(client),
			updatePool(client),
			deletePool(client),
			getNode(client),
			updateNode(client),
		},
	}
}

func get(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get load balancer details",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "list all load balancers for location"},
			cli.StringFlag{Name: "id", Usage: "load balancer id"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" {
				fmt.Printf("location flag is required.\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if c.Bool("all") || c.String("id") == "" {
				resp, err := client.LB.GetAll(c.String("location"))
				if err != nil {
					fmt.Printf("failed to get %s\n", c.Args().First())
					return
				}
				b, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					fmt.Printf("%s\n", err)
					return
				}
				fmt.Printf("%s\n", b)
				return
			}
			resp, err := client.LB.Get(c.String("location"), c.String("id"))
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
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
		Usage:   "create shared load balancer",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "load balancer name [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "description, d", Usage: "load balancer description"},
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			name := c.String("name")
			loc := c.String("location")
			if name == "" || loc == "" {
				fmt.Printf("missing required flags to create load balancer. [use --help to show required flags]\n")
				return
			}

			lb := lb.LoadBalancer{Name: name, Description: c.String("description")}
			resp, err := client.LB.Create(loc, lb)
			if err != nil {
				fmt.Printf("failed to create load balancer [%s] in %s\n", name, loc)
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func update(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "update shared load balancer",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name, n", Usage: "load balancer name [required]"},
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "description, d", Usage: "load balancer description [required]"},
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			id := c.String("id")
			loc := c.String("location")
			name := c.String("name")
			desc := c.String("description")
			if id == "" || loc == "" || name == "" || desc == "" {
				fmt.Printf("missing required flags to update load balancer. [use --help to show required flags]\n")
				return
			}

			lb := lb.LoadBalancer{Name: name, Description: desc}
			err = client.LB.Update(loc, id, lb)
			if err != nil {
				fmt.Printf("failed to update load balancer [%s] in %s\n", lb.Name, loc)
				return
			}
			fmt.Printf("successfully updated lb [%s]\n", id)
		},
	}
}

func delete(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete load balancer details",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("id") == "" || c.String("location") == "" {
				fmt.Printf("missing required flags to delete load balancer. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = client.LB.Delete(c.String("location"), c.String("id"))
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			fmt.Printf("successfully delete lb [%s]\n", c.String("id"))
		},
	}
}
func getPool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get-pool",
		Aliases: []string{"gp"},
		Usage:   "get load balancer pool details",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "list all load balancers for location"},
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool, p", Usage: "load balancer pool id"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" || c.String("id") == "" {
				fmt.Printf("missing required flags to get pool. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if c.Bool("all") || c.String("pool") == "" {
				resp, err := client.LB.GetAllPools(c.String("location"), c.String("id"))
				if err != nil {
					fmt.Printf("failed to get %s\n", c.Args().First())
					return
				}
				b, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					fmt.Printf("%s\n", err)
					return
				}
				fmt.Printf("%s\n", b)
				return
			}
			resp, err := client.LB.GetPool(c.String("location"), c.String("id"), c.String("pool"))
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func createPool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "create-pool",
		Aliases: []string{"cp"},
		Usage:   "create load balancer pool",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.IntFlag{Name: "port", Usage: "pool port [required]"},
			cli.BoolFlag{Name: "sticky", Usage: "use stick persistence"},
			cli.BoolFlag{Name: "standard", Usage: "use standard persistence [default]"},
			cli.BoolFlag{Name: "least-connection, lc", Usage: "use least-connection load balacing"},
			cli.BoolFlag{Name: "round-robin, rr", Usage: "use round-robin load balacing [default]"},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("sticky") && c.Bool("standard") {
				fmt.Println("only one of sticky and standard can be selected")
				return fmt.Errorf("")
			}

			if c.Bool("least-connection") && c.Bool("round-robin") {
				fmt.Println("only one of least-connection and round-robin can be selected")
				return fmt.Errorf("")
			}

			if c.String("id") == "" || c.String("location") == "" || c.Int("port") == 0 {
				fmt.Println("missing required flags, --help for more details")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			pool := lb.Pool{Port: c.Int("port")}
			if c.Bool("sticky") {
				pool.Persistence = lb.Sticky
			} else {
				pool.Persistence = lb.Standard
			}

			if c.Bool("least-connection") {
				pool.Method = lb.LeastConn
			} else {
				pool.Method = lb.RoundRobin
			}

			resp, err := client.LB.CreatePool(c.String("location"), c.String("id"), pool)
			if err != nil {
				fmt.Printf("failed to create load balancer pool for [%s] in %s\n", c.String("id"), c.String("location"))
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func updatePool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "update-pool",
		Aliases: []string{"up"},
		Usage:   "update load balancer pool",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool, p", Usage: "load balancer pool id [required]"},
			cli.BoolFlag{Name: "sticky", Usage: "use stick persistence"},
			cli.BoolFlag{Name: "standard", Usage: "use standard persistence [default]"},
			cli.BoolFlag{Name: "least-connection, lc", Usage: "use least-connection load balacing"},
			cli.BoolFlag{Name: "round-robin, rr", Usage: "use round-robin load balacing [default]"},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("sticky") && c.Bool("standard") {
				fmt.Println("only one of sticky and standard can be selected")
				return fmt.Errorf("")
			}

			if c.Bool("least-connection") && c.Bool("round-robin") {
				fmt.Println("only one of least-connection and round-robin can be selected")
				return fmt.Errorf("")
			}

			if c.String("id") == "" || c.String("location") == "" || c.String("pool") == "" {
				fmt.Println("missing required flags, --help for more details")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			pool := lb.Pool{}
			if c.Bool("sticky") {
				pool.Persistence = lb.Sticky
			} else {
				pool.Persistence = lb.Standard
			}

			if c.Bool("least-connection") {
				pool.Method = lb.LeastConn
			} else {
				pool.Method = lb.RoundRobin
			}

			err = client.LB.UpdatePool(c.String("location"), c.String("id"), c.String("pool"), pool)
			if err != nil {
				fmt.Printf("failed to update load balancer pool [%s] in %s\n", c.String("pool"), c.String("location"))
				return
			}
			fmt.Printf("successfully updated pool [%s]\n", c.String("pool"))
		},
	}
}

func deletePool(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "delete-pool",
		Aliases: []string{"dp"},
		Usage:   "delete load balancer pool",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool, p", Usage: "load balancer pool id [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" || c.String("id") == "" || c.String("pool") == "" {
				fmt.Printf("missing required flags to get pool. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = client.LB.DeletePool(c.String("location"), c.String("id"), c.String("pool"))
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			fmt.Printf("sucessfully deleted pool [%s] from lb [%s]\n", c.String("pool"), c.String("id"))
		},
	}
}

func getNode(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "get-node",
		Aliases: []string{"gn"},
		Usage:   "get load balancer node details",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool, p", Usage: "load balancer pool id [required]"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" || c.String("id") == "" || c.String("pool") == "" {
				fmt.Printf("missing required flags to get pool. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resp, err := client.LB.GetAllNodes(c.String("location"), c.String("id"), c.String("pool"))
			if err != nil {
				fmt.Printf("failed to get %s\n", c.Args().First())
				return
			}
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("%s\n", b)
		},
	}
}

func updateNode(client *clc.Client) cli.Command {
	return cli.Command{
		Name:    "update-node",
		Aliases: []string{"un"},
		Usage:   "update load balancer node details",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "id", Usage: "load balancer id [required]"},
			cli.StringFlag{Name: "location, l", Usage: "load balancer location [required]"},
			cli.StringFlag{Name: "pool, p", Usage: "load balancer pool id [required]"},
			cli.StringSliceFlag{Name: "host, h", Usage: "node hostname and port (ex. 10.10.10.10:8080)"},
		},
		Before: func(c *cli.Context) error {
			if c.String("location") == "" || c.String("id") == "" || c.String("pool") == "" {
				fmt.Printf("missing required flags to get pool. [use --help to show required flags]\n")
				return fmt.Errorf("")
			}
			return nil
		},
		Action: func(c *cli.Context) {
			client, err := util.MaybeLoadConfig(c, client)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			nodes := make([]lb.Node, len(c.StringSlice("host")))
			for i, v := range c.StringSlice("host") {
				split := strings.Split(v, ":")
				port, err := strconv.Atoi(split[1])
				if err != nil {
					fmt.Printf("failed parsing %s", v)
					return
				}
				node := lb.Node{
					IPaddress:   split[0],
					PrivatePort: port,
				}
				nodes[i] = node
			}
			err = client.LB.UpdateNodes(c.String("location"), c.String("id"), c.String("pool"), nodes...)
			if err != nil {
				fmt.Printf("failed to update nodes\n")
				return
			}
			fmt.Printf("nodes updates for pool %s\n", c.String("pool"))
		},
	}
}
