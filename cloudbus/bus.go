package cloudbus

import (
	"encoding/json"
	"errors"
	"microlight/aggregate"
	"microlight/broker"
	"microlight/cloudstream"
	"microlight/logger"
	"microlight/message"
)

var (
	log = logger.Create()
)

//bus
type bus struct {
	subscribers   map[string]*dispatcher
	EventSource   chan *message.Event
	stop          chan int
	streamBuilder cloudstream.StreamBuilder
	stream        cloudstream.Stream
}

type dispatcher struct {
	builder   aggregate.AggregateBuilder
	instances map[string]aggregate.Aggregate
}

//Bus ...
type Bus interface {
	SendCommand(command *message.Command) error
	RegisterAggregate(AggregateName string, builder aggregate.AggregateBuilder)
	Start()
	Stop()
}

//CreateBus ...
func CreateBus(ClusterID string, ListenerID string) Bus {
	st := cloudstream.CreateStreamBuilder(broker.CreateNatsBroker(ClusterID, ListenerID))
	b := &bus{
		streamBuilder: st,
		subscribers:   make(map[string]*dispatcher),
		EventSource:   make(chan *message.Event, 10),
	}

	return b
}

//Start ...
func (b *bus) Start() {
	log.Info.Println("Start")

	for key := range b.subscribers {
		b.streamBuilder.RegisterSink(key, key,
			func(msg interface{}) error {
				var HandlinError error
				log.Debug.Println("I recive a command")
				v := message.Command{}
				json.Unmarshal(msg.([]byte), &v)
				subscriber := b.subscribers[v.AggregateType]
				if subscriber == nil {
					HandlinError = errors.New("aggregate type is not registered")
					return HandlinError
				}
				aggID := v.AggregateID
				if aggID == nil {
					aggID = ""
				}
				aggregate := subscriber.instances[aggID.(string)]
				if aggregate == nil {
					aggregate = subscriber.builder.Initialize(v.AggregateID)
					go func() {
						for event := range aggregate.GetEventSource() {

							b.EventSource <- event
						}
					}()
					if id := aggregate.GetID(); id != "" {
						subscriber.instances[id] = aggregate
					}
				}
				aggregate.Command(&v)

				return HandlinError

			})
	}
	b.stream = b.streamBuilder.Build()
	go func() {
		for event := range b.EventSource {
			b.stream.Publish(event.AggregateType+"-events", event)
		}
	}()
	log.Debug.Println("waiting")
	//stat := <-b.stop

}

//Stop ...
func (b *bus) Stop() {
	log.Info.Println("stopping")
	b.stream.Close()
	log.Info.Println("Shotdown")
	b.stop <- 0
}

//SendCommand ...
func (b *bus) SendCommand(command *message.Command) error {
	return b.stream.Publish(command.AggregateType, command)
}

//RegisterAggregate ...
func (b *bus) RegisterAggregate(AggregateName string, builder aggregate.AggregateBuilder) {
	b.subscribers[AggregateName] = &dispatcher{builder: builder,
		instances: make(map[string]aggregate.Aggregate)}
}
