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

	switch event.EventType {
	case "state_changed":
		if err := m.handleStateChanged(event); err != nil {
			return err
		}
	}

	return m.nc.Publish(fmt.Sprintf("hass.%s.events", m.Identifier()), b)
}

func (m *Module) handleStateChanged(event go_hass_ws.Event) error {
	var evt go_hass_ws.StateChangedEvent
	if err := json.Unmarshal(*event.Data, &evt); err != nil {
		return fmt.Errorf("unable to decode state_changed_event: %w", err)
	}

	if _, err := m.kvEntities.Put(evt.EntityId, evt.NewState); err != nil {
		return fmt.Errorf("unable to persist state of entity %s: %w", evt.EntityId, err)
	}

	return nil
}

func (m *Module) handleEntityRegistryUpdated(event go_hass_ws.Event) error {
	var evt EntityRegistryUpdated
	if err := json.Unmarshal(*event.Data, &evt); err != nil {
		return fmt.Errorf("unable to decode entity_registry_updated: %w", err)
	}

	switch evt.Action {
	case "create":
		if _, err := m.kvEntities.Put(evt.EntityId, []byte("{}")); err != nil {
			return fmt.Errorf("unable to persist state of entity %s: %w", evt.EntityId, err)
		}
	case "update":
		log.Warn().Msg("Not implemented yet.")
	case "remove":
		if err := m.kvEntities.Delete(evt.EntityId); err != nil {
			return fmt.Errorf("unable to remove entity %s: %w", evt.EntityId, err)
		}
	}

	return nil
}

func (m *Module) handleCallService(event go_hass_ws.Event) error {
	return nil
}

type EntityRegistryUpdated struct {
	Action   string          `json:"action"`
	EntityId string          `json:"entity_id"`
	Changes  json.RawMessage `json:"changes,omitempty"`
}
