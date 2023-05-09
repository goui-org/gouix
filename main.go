package main

import (
	"fmt"
	"log"
	"os"

	"main/build"
	"main/create"
	"main/serve"

	"github.com/urfave/cli/v2"
)

func main() {
	os.Setenv("GOOS", "js")
	os.Setenv("GOARCH", "wasm")
	app := &cli.App{
		Name:  "goui",
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
					args := c.Args()
					fmt.Println(args)
					return create.Create(args.First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
