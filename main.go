package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-inject"
	app.Usage = "Copy files/directories from hosts to running Docker containers"
	app.Version = "0.0.0"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.HelpFlag,
	}
	cli.AppHelpTemplate = appHelpTemplate
	app.Action = func(c *cli.Context) {
		inj, err := newInjector(os.Stderr, c.Args())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", app.Name, err)
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		if err := inj.run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", app.Name, err)
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
	}
	app.Run(os.Args)
}
