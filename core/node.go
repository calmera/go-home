package core

import (
	"fmt"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	log "github.com/rs/zerolog/log"
	"sync"
	"time"
)

type NodeConfig struct {
	StorageDir string `json:"storage_dir"`
	ConfigDir  string `json:"config_dir"`
}

func NewNode(cfg NodeConfig) (*Node, error) {
	n := &Node{
		config:  cfg,
		ns:      nil,
		nc:      nil,
		Modules: map[string]Module{},
		wg:      &sync.WaitGroup{},
	}

	return n, nil
}

type Node struct {
	config  NodeConfig
	ns      *server.Server
	nc      *nats.Conn
	Modules map[string]Module
	wg      *sync.WaitGroup
}

func (n *Node) RegisterModule(module Module) {
	n.Modules[module.Identifier()] = module
}

func (n *Node) Start() error {
	opts := &server.Options{
		JetStream: true,
		StoreDir:  n.config.StorageDir,
		MQTT:      server.MQTTOpts{},
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("unable to create the embedded nats server: %w", err)
	}

	n.ns = ns

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		n.Close()
		return fmt.Errorf("not ready for connection")
	}

	log.Info().Msg("Embedded NATS server started")

	// -- connect to the embedded nats server
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		n.Close()
		return err
	}

	n.nc = nc

	for _, m := range n.Modules {
		log.Info().Str("module", m.Identifier()).Msg("starting module")
		go func(mod Module, wg *sync.WaitGroup, nc *nats.Conn) {
			wg.Add(1)
			err := mod.Run(nc)
			if err != nil {
				log.Err(err)
			} else {
				log.Info().Str("module", mod.Identifier()).Msg("module finished")
			}
			wg.Add(-1)
		}(m, n.wg, nc)
	}

	log.Info().Msg("Node Ready!")

	return nil
}

func (n *Node) Wait() {
	n.ns.WaitForShutdown()
}

func (n *Node) Close() {
	// -- signal all modules to close
	for _, m := range n.Modules {
		m.Close()
	}

	// -- wait for all modules to be closed
	n.wg.Wait()

	// -- close the client
	n.nc.Close()

	// -- close the server
	n.ns.Shutdown()
}
