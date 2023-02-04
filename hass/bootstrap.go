package hass

import (
	"encoding/json"
	"fmt"
	go_hass_ws "github.com/calmera/go-hass-ws"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func (m *Module) bootstrap() error {
	if err := m.bootstrapServices(); err != nil {
		return fmt.Errorf("unable to bootstrap services: %w", err)
	}

	if err := m.bootstrapEntities(); err != nil {
		return fmt.Errorf("unable to bootstrap entities: %w", err)
	}

	return nil
}

func (m *Module) bootstrapServices() error {
	bucket := fmt.Sprintf("hass_%s_services", m.Identifier())

	kv, err := m.js.KeyValue(bucket)
	if err != nil {
		if err == nats.ErrBucketNotFound {
			kv2, err2 := m.js.CreateKeyValue(&nats.KeyValueConfig{
				Bucket:  bucket,
				Storage: nats.FileStorage,
			})
			if err2 != nil {
				return err2
			}

			kv = kv2
		} else {
			return err
		}
	}

	err = m.hass.GetServices(func(services map[string]go_hass_ws.ServiceDomain) {
		for sid, service := range services {
			b, err := json.Marshal(service)
			if err != nil {
				log.Err(fmt.Errorf("unable to encode hass service to json: %w", err))
			}

			if _, err := kv.Create(sid, b); err != nil {
				log.Err(fmt.Errorf("unable to store hass service %s to json: %w", sid, err))
			}
		}
	})
	if err != nil {
		return err
	}

	m.kvServices = kv
	log.Info().Str("module", m.Identifier()).Msg("services registered")

	return nil
}

func (m *Module) bootstrapEntities() error {
	bucket := fmt.Sprintf("hass_%s_entities", m.Identifier())

	kv, err := m.js.KeyValue(bucket)
	if err != nil {
		if err == nats.ErrBucketNotFound {
			kv2, err2 := m.js.CreateKeyValue(&nats.KeyValueConfig{
				Bucket:  bucket,
				Storage: nats.FileStorage,
			})
			if err2 != nil {
				return err2
			}

			kv = kv2

		} else {
			return err
		}
	}

	err = m.hass.GetStates(func(states map[string]go_hass_ws.State) {
		for eid, entity := range states {
			b, err := json.Marshal(entity)
			if err != nil {
				log.Err(fmt.Errorf("unable to encode hass state to json: %w", err))
			}

			if _, err := kv.Create(eid, b); err != nil {
				log.Err(fmt.Errorf("unable to store hass entity %s to json: %w", eid, err))
			}
		}
	})
	if err != nil {
		return err
	}

	m.kvEntities = kv
	log.Info().Str("module", m.Identifier()).Msg("entities registered")

	return nil
}
