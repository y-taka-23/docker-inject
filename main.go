package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-inject"
	app.Usage = "Copy files/directories from hosts to running Docker containers"
	app.Version = "0.0.0"
	app.Action = func(c *cli.Context) {
		newInjector().run(c.Args())
	}
	app.Run(os.Args)
}
