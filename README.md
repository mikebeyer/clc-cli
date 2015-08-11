clc cli
===================
Insallation
---------------------

```sh
$ go get github.com/mikebeyer/clc-cli
$ make deps
$ make build
```

Configuration
------
In order to use the clc cli tool you need a configuration.

It does this by first looking into the environment, and then falling back to a well known (./config.json) config location.

The cli tool can create a template for a configuration file with `./clc --gen-config`

**Env Vars**

```sh
CLC_USERNAME
CLC_PASSWORD
CLC_ALIAS
```

**Config File**

```json
{
  "user": {
    "username": "USERNAME",
    "password": "PASSWORD"
  },
  "alias": "DEFAULT-ALIAS"
}
```


Usage
------
```sh
NAME:
   clc - clc v2 api cli

USAGE:
   clc [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR(S):
   Mike Beyer <michael.beyer@ctl.io>

COMMANDS:
   server, s			server api
   status				status api
   anti-alias, aa		anti-alias api
   alert, a				alert api
   load-balancer, lb	load balancer api
   group, g				group api
   help, h				Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gen-config			create template configuration file
   --help, -h			show help
   --version, -v		print the version

```

Example
------

```sh
./clc server create -n 'server' -c 1 -m 1 -g GROUP-ID -t standard --ubuntu-14
```

```sh
./clc server ssh SERVERNAME
```

License
-------
This project is licensed under the [Apache License v2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
