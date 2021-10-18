package broker

type Broker interface {
	Publish(topic string, message interface{}) error

	Subscribe(clientId string, topic string, handler func(message interface{}) error) error

	Unsubscribe(clientId string) error
}
