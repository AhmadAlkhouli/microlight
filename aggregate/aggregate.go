package aggregate

import (
	"microlight/message"
)

type Aggregate interface {
	EventsAggregate()

	GetCommandHandler(key string) func(Aggregate, message.Command)

	GetID() string

	Command(*message.Command)

	GetEventSource() chan *message.Event

	RaiseEvent(message.Event)

	GetEntity() interface{}
}
