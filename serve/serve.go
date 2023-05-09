package serve

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/twharmon/gouix/devserver"
)

func Start() error {
	os.Setenv("DEBUG", "true")
	port := os.Getenv("PORT")
	if port == "" {
		os.Setenv("PORT", "3000")
	}
	server, err := devserver.New()
	if err != nil {
		return fmt.Errorf("devserver.New: %w", err)
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
