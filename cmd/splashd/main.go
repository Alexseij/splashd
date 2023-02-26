package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/splashd/config"
	"github.com/splashd/internal/core"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "splashd",
		Flags: SetupFlags(),
		Action: func(ctx *cli.Context) error {
			daemonLog := logrus.WithField(
				"Daemon",
				"splashd",
			)

			d := core.NewDaemon(
				config.DaemonOpts.Socket,
				daemonLog,
			)

			return d.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "socket",
			Aliases:     []string{"s"},
			Destination: &config.DaemonOpts.Socket,
			Value:       core.DefaultSocketValue,
		},
	}
}
