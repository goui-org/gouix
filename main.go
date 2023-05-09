package main

import (
	"log"
	"os"

	"github.com/twharmon/gouix/build"
	"github.com/twharmon/gouix/create"
	"github.com/twharmon/gouix/serve"

	"github.com/urfave/cli/v2"
)

func main() {
	os.Setenv("GOOS", "js")
	os.Setenv("GOARCH", "wasm")
	app := &cli.App{
		Name:  "gouix",
		Usage: "develop applications with goui",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "start develpoment server",
				Action: func(c *cli.Context) error {
					return serve.Start()
				},
			},
			{
				Name:  "build",
				Usage: "build application",
				Action: func(c *cli.Context) error {
					return build.New().Run()
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
