package hass

import (
	go_hass_ws "github.com/calmera/go-hass-ws"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"sync"
)

func NewModule(config ModuleConfig) (*Module, error) {
	return &Module{config: config}, nil
}

type Module struct {
	config ModuleConfig
	wg     sync.WaitGroup
	nc     *nats.Conn
	js     nats.JetStreamContext
	hass   *go_hass_ws.HassClient

	kvServices nats.KeyValue
	kvEntities nats.KeyValue
}

func (m *Module) Identifier() string {
	return m.config.Id
}

func (m *Module) Run(nc *nats.Conn, js nats.JetStreamContext) error {
	m.nc = nc
	m.js = js

	hass, err := go_hass_ws.Connect(go_hass_ws.Config{
		Url:   m.config.Url,
		Token: m.config.Token,
	})
	if err != nil {
		panic(err)
	}
	defer hass.Close()
	m.hass = hass

	// -- bootstrap
	if err := m.bootstrap(); err != nil {
		return err
	}

	// -- start listening for events
	if _, err := m.hass.Subscribe(func(msg go_hass_ws.HassAuthenticatedMessage) {}, m.handleEvent); err != nil {
		return err
	}

	log.Info().Str("module", m.Identifier()).Msg("listening for HASS events")

	// -- wait until our hass connection dies. might not be the best idea, but will be good for the time being
	hass.WaitUntilAllHandled()
	return nil
}

func (m *Module) Close() {
	m.wg.Done()
}
