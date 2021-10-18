package test

import (
	"microlight/eventstore"
	"microlight/logger"
	"microlight/message"
)

type eventStoreMock struct{}

var log = logger.Create()

func NewEventStoreMock() eventstore.EventStore {
	return &eventStoreMock{}
}
func (e *eventStoreMock) Persist(event interface{}) {
	log.Debug.Printf("push the event %+v", event.(*message.Event))
}
