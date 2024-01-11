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
		Usage: "develop user interfaces with goui",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "tinygo",
				Usage: "use the tinygo compiler",
			},
			&cli.StringFlag{
				Name:  "proxy",
				Usage: "proxy requests",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "start develpoment server",
				Action: func(c *cli.Context) error {
					return serve.Start(c.Bool("tinygo"), c.String("proxy"))
				},
			},
			{
				Name:  "build",
				Usage: "build application",
				Action: func(c *cli.Context) error {
					return build.New(c.Bool("tinygo")).Run()
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
