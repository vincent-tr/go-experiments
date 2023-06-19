package store

import "mylife-tools-server/services/io/serialization"

type Entity interface {
	serialization.Marshaller
	Id() string
}
