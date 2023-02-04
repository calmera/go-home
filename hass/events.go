package hass

import (
	"encoding/json"
	"fmt"
	go_hass_ws "github.com/calmera/go-hass-ws"
	"github.com/rs/zerolog/log"
)

func (m *Module) handleEvent(id uint64, event go_hass_ws.Event) error {
	b, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	log.Debug().
		Str("module", m.Identifier()).
		Str("eventType", event.EventType).Msg(string(b))

	return m.nc.Publish(fmt.Sprintf("hass.%s.events", m.Identifier()), b)
}
