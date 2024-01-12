package serve

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/twharmon/gouix/config"
	"github.com/twharmon/gouix/server"
)

func Start(cfg *config.Config) error {
	os.Setenv("DEBUG", "true")
	server, err := server.New(cfg)
	if err != nil {
		return fmt.Errorf("serve.Start: %w", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		server.Shutdown()
		os.Exit(0)
	}()
	return server.Run()
}
