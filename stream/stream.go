package stream

import (
	"microlight/broker"
	"microlight/logger"
)

var (
	log = logger.Create()
)

type consumer func(message interface{}) error

type producer func(stop chan interface{})

type (
	Sink struct {
		sink  consumer
		topic string
	}
	Source struct {
		source producer
		topic  string
	}
	stream struct {
		binder  broker.Broker
		sinks   map[string]*Sink
		sources map[string]*Source
		stop    chan int
	}
	StreamBuilder interface {
		Build() Stream
		RegisterSource(id string, topic string, source producer) StreamBuilder
		RegisterSink(id string, topic string, sink consumer) StreamBuilder
	}
	Stream interface {
		Publish(topic string, msg interface{}) error
		DiscontinueSource(id string)
		DiscontinueSink(id string)
		Close()
	}
)

func CreateStreamBuilder(broker broker.Broker) StreamBuilder {
	return &stream{binder: broker, sources: make(map[string]*Source), sinks: make(map[string]*Sink)}
}

func (s *stream) Build() Stream {
	go func() {
		s.startListening()
		s.startDataGenerator()
		log.Debug.Println("build ----------------")

	}()

	log.Debug.Println("end ----------------")

	return s
}
func (s *stream) RegisterSource(id string, topic string, source producer) StreamBuilder {

	s.sources[id] = &Source{topic: topic,
		source: source,
	}
	return s
}
func (s *stream) RegisterSink(id string, topic string, sink consumer) StreamBuilder {
	s.sinks[id] = &Sink{topic: topic, sink: sink}
	return s
}
func (s *stream) Publish(topic string, msg interface{}) error {
	return s.binder.Publish(topic, msg)
}

func (s *stream) startListening() {

	for key, value := range s.sinks {
		go s.binder.Subscribe(key, value.topic, value.sink)
		log.Debug.Printf("Listening %s to topic %s", key, value.topic)

	}
}
func (str *stream) startDataGenerator() {
	for _, value := range str.sources {

		go func(s *Source) {
			channel := make(chan interface{})
			go s.source(channel)

			log.Debug.Println("start generating")

			for msg := range channel {
				log.Debug.Println("generating")

				str.Publish(s.topic, msg)
			}
			log.Debug.Println("end generating")

		}(value)
	}
}
func (s *stream) DiscontinueSource(id string) {
	//close(stream.sources[id].channel)
}
func (s *stream) DiscontinueSink(id string) {
	s.binder.Unsubscribe(id)

}
func (s *stream) Close() {
	for i := range s.sinks {
		s.binder.Unsubscribe(i)
	}
	/* for _ = range s.sources {
		//close(stream.sources[j].channel)
	} */
	s.binder = nil
	s.stop <- 0
}
