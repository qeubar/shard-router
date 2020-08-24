package main

import (
	"fmt"
	"os"

	"shard-router/api"

	"github.com/urfave/cli/v2"
)

var appCommands = []*cli.Command{
	{
		Name:   "server",
		Usage:  "start the API HTTP server",
		Action: api.StartServer,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "proxy-protocol",
				Usage: "Expect the proxy protocol",
			},
		},
	},
}

func main() {
	app := &cli.App{
		Name:     "shard-router",
		Usage:    "Shard-aware web routing and cli",
		Version:  "0.0.1",
		Commands: appCommands,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
