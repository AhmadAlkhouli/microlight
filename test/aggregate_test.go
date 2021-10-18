package test

import (
	"microlight/aggregate"
	"microlight/message"
	"testing"
	"time"
)

type User struct {
	ID   string
	name string
}
type NewCommand struct {
	UserID string `json:"UserID"`
	Name   string `json:"Name"`
}
type UpdateCommand struct {
	userID string
	name   string
}
type UserCreated struct {
	UserID string `json:"UserID"`
	Name   string `json:"Name"`
}
type UserUpdated struct {
	userID string
	name   string
}

var (
	//log               = logger.Create()
	NewCommandHandler = func(aggr aggregate.Aggregate, c message.Command) {
		com := NewCommand{}

		value, ok := c.CommandData.(map[string]interface{})
		if ok {
			message.FillStruct(&com, value)
		} else {
			com = c.CommandData.(NewCommand)
		}

		data := UserCreated{UserID: com.UserID, Name: com.Name}
		log.Debug.Println("rais events")
		aggr.RaiseEvent(message.Event{EventType: "UserCreated", AggregateID: "1", AggregateType: "User", EventData: data})
		aggr.EventsAggregate()

	}
	UpdateCommandHandler = func(aggr aggregate.Aggregate, c message.Command) {
		command := c.CommandData.(UpdateCommand)
		data := UserUpdated{userID: command.userID, name: command.name}
		log.Debug.Println("rais events")
		aggr.RaiseEvent(message.Event{EventType: "UserUpdated", AggregateID: "1", AggregateType: "User", EventData: &data})
		aggr.EventsAggregate()

	}
	UserCreatedHandler = func(i interface{}, e message.Event) {
		ent := i.(*User)
		ent.name = e.EventData.(UserCreated).Name

	}
	UserUpdatedHandler = func(i interface{}, e message.Event) {
		ent := i.(*User)
		ent.name = e.EventData.(*UserUpdated).name

	}
	PersistUserData = func(i interface{}) {
		log.Debug.Println("end")

		user := i.(*User)
		log.Debug.Printf("id=%s, name=%s", user.ID, user.name)

	}
	BeginState = func(ID interface{}) interface{} {
		return &User{ID: "1", name: "Ahmad"}
	}
)

func TestAggregate(t *testing.T) {
	t.Log("start")
	builder := aggregate.NewAggregateBuilder(NewEventStoreMock(), "")

	aggr := builder.Initializer(BeginState).
		Handle("new", NewCommandHandler).
		Handle("update", UpdateCommandHandler).
		ApplyEvent("UserCreated", UserCreatedHandler).
		ApplyEvent("UserUpdated", UserUpdatedHandler).
		Persistent(PersistUserData).
		Initialize("1")

	aggr.Command(message.CreateCommand(
		"new", "1", "User", NewCommand{UserID: "1", Name: "Ahmad Al Khouli"}))

	aggr.Command(message.CreateCommand(
		"update", "1", "User", UpdateCommand{userID: "1", name: "Ahmad Al Khouli1"}))

	for i := range aggr.GetEventSource() {
		println(i.AggregateID.(string))
	}
	time.Sleep(5 * time.Second)

}
