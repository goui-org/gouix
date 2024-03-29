package main

import (
	"log"
	"os"

	"github.com/goui-org/gouix/build"
	"github.com/goui-org/gouix/config"
	"github.com/goui-org/gouix/create"
	"github.com/goui-org/gouix/serve"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gouix",
		Usage: "develop user interfaces with goui",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "start develpoment server",
				Action: func(c *cli.Context) error {
					return serve.Start(config.Get())
				},
			},
			{
				Name:  "build",
				Usage: "build application",
				Action: func(c *cli.Context) error {
					return build.New(config.Get()).Run()
				},
			},
			{
				Name:  "create",
				Usage: "create a new goui application",
				Action: func(c *cli.Context) error {
					return create.Create(c.Args().First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
