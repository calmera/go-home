package hass

import "github.com/nats-io/nats.go"

type HassModule struct {
}

func (m *HassModule) Identifier() string {
	return "hass"
}

func (m *HassModule) Run(nc *nats.Conn) error {
	return nil
}

func (m *HassModule) Close() {

}
