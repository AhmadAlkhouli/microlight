package aggregate

import (
	"microlight/message"
)

// AggregateBuilder ...
type AggregateBuilder interface {

	// Consume the event (apply side effect into the aggregate model)
	ApplyEvent(eventName string, concumer func(interface{}, message.Event)) AggregateBuilder

	// Define Handler for specific command that will be triggered when recive command message
	// the handler will rais domain events as side effects
	Handle(commandName string, handler func(Aggregate, message.Command)) AggregateBuilder

	// Define persistance function to save the updated domains after aggreate
	Persistent(func(interface{})) AggregateBuilder

	Initializer(func(ID interface{}) interface{}) AggregateBuilder

	Initialize(ID interface{}) Aggregate
}
