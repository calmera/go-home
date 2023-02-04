package core

import (
	"github.com/nats-io/nats.go"
)

type Module interface {
	Identifier() string
	Run(nc *nats.Conn, js nats.JetStreamContext) error
	Close()
}
