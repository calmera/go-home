package main

import (
	"fmt"
	"github.com/calmera/go-home/core"
	"github.com/calmera/go-home/hass"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := core.NodeConfig{
		StorageDir: "data",
		ConfigDir:  "config",
	}

	hassModules, err := hass.LoadHassConfig(fmt.Sprintf("%s/hass.json", cfg.ConfigDir))
	if err != nil {
		panic(err)
	}

	n, err := core.NewNode(cfg)
	if err != nil {
		panic(err)
	}

	for _, hm := range hassModules {
		n.RegisterModule(hm)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(n *core.Node) {
		<-c
		n.Close()
		os.Exit(1)
	}(n)

	if err := n.Start(); err != nil {
		panic(err)
	}

	n.Wait()
}
