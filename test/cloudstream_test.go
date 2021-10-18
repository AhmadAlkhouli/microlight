package test

import (
	"fmt"
	"microlight/stream"
	"microlight/transport"
	"testing"
	"time"
)

var (
	testStream = stream.CreateStreamBuilder(transport.CreateNatsBroker("test-cluster", "test"))
)

func TestStream(t *testing.T) {
	testStream.RegisterSource("s1", "s1", func(channel chan interface{}) {
		defer close(channel)
		log.Debug.Println("hello")
		for i := 0; i < 5; i++ {
			log.Debug.Printf("generate value %d", i)
			value := fmt.Sprintf("hello %d", i)
			channel <- value
		}
	}).
		RegisterSink("sink1", "s1", func(message interface{}) error {
			log.Debug.Printf("%s ", message)
			return nil
		}).Build()
	log.Debug.Println("end===================")

	time.Sleep(5 * time.Second)
}
