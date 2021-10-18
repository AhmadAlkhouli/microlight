package aggregate

import (
	"microlight/eventstore"
	"microlight/logger"
	"microlight/message"
	"reflect"
	"sync"
)

type (
	baseAggregate struct {
		ID               string
		Entity           interface{}
		events           []*message.Event
		eventsConcumers  map[string]func(interface{}, message.Event)
		commandsHandlers map[string]func(Aggregate, message.Command)
		persistent       func(interface{})
		initializer      func(ID interface{}) interface{}
		eventStore       eventstore.EventStore
		commandSink      chan *message.Command
		eventSource      chan *message.Event
		state            int
		lock             sync.RWMutex
	}
)

var (
	log = logger.Create()
)

//NewAggregateBuilder ...
func NewAggregateBuilder(eventStore eventstore.EventStore, aggregateID string) AggregateBuilder {
	aggr := &baseAggregate{
		events:           make([]*message.Event, 0),
		eventsConcumers:  make(map[string]func(interface{}, message.Event)),
		commandsHandlers: make(map[string]func(Aggregate, message.Command)),
		eventStore:       eventStore,
		ID:               aggregateID,
		commandSink:      make(chan *message.Command, 10),
		eventSource:      make(chan *message.Event, 10),
		lock:             sync.RWMutex{},
		state:            1,
	}

	return aggr
}

//Initialize ...
func (aggr *baseAggregate) Initialize(ID interface{}) Aggregate {

	aggr.Entity = aggr.initializer(ID)
	if ID != nil {
		aggr.ID = ID.(string)
	}
	go func() {

		for {

			select {
			case com, ok := <-aggr.commandSink:
				if !ok {
					log.Error.Println("sink is closed")
					return
				}
				log.Debug.Println("Start Aggregate")
				log.Debug.Println(",," + com.AggregateType + ",," + com.CommandType)
				log.Debug.Printf("%+v\n", com)
				handler := aggr.GetCommandHandler(com.CommandType)
				handler(aggr, (*com))
				log.Info.Println("Handled")
			default:
				if aggr.state == 0 {
					log.Info.Println("stop")
					close(aggr.commandSink)
					close(aggr.eventSource)
					aggr = nil
					return
				}
			}
		}

	}()
	log.Info.Println("exit init routine")
	return aggr
}

//RaiseEvent ...
func (aggr *baseAggregate) RaiseEvent(event message.Event) {
	aggr.lock.Lock()
	aggr.events = append(aggr.events, &event)
	aggr.lock.Unlock()
}

//EventsAggregate ...
func (aggr *baseAggregate) EventsAggregate() {
	defer func() {
		aggr.events = make([]*message.Event, 0)
	}()
	aggr.lock.Lock()
	log.Debug.Printf("events count is %d", len(aggr.events))
	for _, element := range aggr.events {
		handler := aggr.eventsConcumers[element.EventType]
		if handler != nil {
			handler(aggr.Entity, *element)
		}
		aggr.eventStore.Persist(element)
		log.Debug.Println("Persist")
		aggr.eventSource <- element

	}
	//Prevent Modifying Entity
	entity := cloneValue(aggr.Entity)
	aggr.persistent(entity)
	aggr.state = 0
	aggr.lock.Unlock()

}

//ApplyEvent ...
func (aggr *baseAggregate) ApplyEvent(eventName string,
	concumer func(interface{}, message.Event)) AggregateBuilder {
	aggr.eventsConcumers[eventName] = concumer
	return aggr
}

//Handle ...
func (aggr *baseAggregate) Handle(commandName string,
	handler func(Aggregate, message.Command)) AggregateBuilder {
	aggr.commandsHandlers[commandName] = handler
	return aggr
}

//GetEntity ...
func (aggr *baseAggregate) GetEntity() interface{} {
	e := aggr.Entity
	return e
}

//Persistent ...
func (aggr *baseAggregate) Persistent(persistent func(interface{})) AggregateBuilder {
	aggr.persistent = persistent
	return aggr
}
func cloneValue(source interface{}) interface{} {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		return y.Interface()
	} else {
		return x.Interface()
	}

}

//Initialize ...
func (aggr *baseAggregate) Initializer(defaultStatus func(Id interface{}) interface{}) AggregateBuilder {
	aggr.initializer = defaultStatus
	return aggr
}

//GetCommandHandler ...
func (aggr *baseAggregate) GetCommandHandler(key string) func(Aggregate, message.Command) {
	handler, ok := aggr.commandsHandlers[key]
	if ok {
		return handler
	}
	return nil
}

func (aggr *baseAggregate) GetID() string {
	return aggr.ID
}

func (aggr *baseAggregate) Command(command *message.Command) {
	aggr.commandSink <- command
}

func (aggr *baseAggregate) GetEventSource() chan *message.Event {
	return aggr.eventSource
}
