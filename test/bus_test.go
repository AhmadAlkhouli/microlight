package test

import (
	"microlight/aggregate"
	"microlight/bus"
	"microlight/logger"
	"microlight/message"
	"testing"
	"time"
)

var (
	logg = logger.Create()
)

func TestBus(t *testing.T) {
	//c := make(chan int)
	builder := aggregate.NewAggregateBuilder(NewEventStoreMock(), "")

	builder.Initializer(BeginState).
		Handle("new", NewCommandHandler).
		Handle("update", UpdateCommandHandler).
		ApplyEvent("UserCreated", UserCreatedHandler).
		ApplyEvent("UserUpdated", UserUpdatedHandler).
		Persistent(PersistUserData)
	bus := bus.CreateBus("test-cluster", "test1")
	bus.RegisterAggregate("User", builder)
	bus.Start()
	logg.Debug.Println("strarted")
	com := message.CreateCommand(
		"new", "1", "User", NewCommand{UserID: "1", Name: "Ahmad Al Khouli"})
	bus.SendCommand(com)
	time.Sleep(5 * time.Second)

}
