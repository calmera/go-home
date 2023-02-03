package main

import (
	"github.com/calmera/go-home/core"
	"github.com/calmera/go-home/hass"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := core.NodeConfig{
		StorageDir: "data",
	}

	n, err := core.NewNode(cfg)
	if err != nil {
		log.Err(err)
	}

	n.RegisterModule(&hass.HassModule{})

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(n *core.Node) {
		<-c
		n.Close()
		os.Exit(1)
	}(n)

	if err := n.Start(); err != nil {
		log.Err(err)
	}

	n.Wait()
}
