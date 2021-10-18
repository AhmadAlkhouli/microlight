package transport

import (
	"encoding/json"
	"microlight/broker"
	"microlight/logger"
	"time"

	"github.com/nats-io/stan.go"
)

var (
	log = logger.Create()
)

// natsBroker ...
type natsBroker struct {
	connection    stan.Conn
	ClusterID     string
	ClientID      string
	subscribtions map[string]func() error
}

// CreateNatsBroker ...
func CreateNatsBroker(clusterID string, clientID string) broker.Broker {
	broker := &natsBroker{
		ClusterID:     clusterID,
		ClientID:      clientID,
		subscribtions: make(map[string]func() error),
	}
	broker.GetConnection()
	return broker
}

// GetConnection ...
func (b *natsBroker) GetConnection() (*stan.Conn, error) {
	if b.connection != nil {
		return &b.connection, nil
	}
	nc, err := stan.Connect(b.ClusterID, b.ClientID, stan.ConnectWait(time.Minute))
	if err != nil {
		log.Info.Panic("No connection")
	}
	b.connection = nc

	return &b.connection, err
}

// Publish ...
func (b *natsBroker) Publish(topic string, message interface{}) error {
	conn, err := b.GetConnection()
	if err != nil {
		return err
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = (*conn).Publish(topic, data)
	log.Debug.Printf("%s send to broker with topic %s", data, topic)
	if err != nil {
		return err
	}

	return nil
}
func (b *natsBroker) Subscribe(clientId string, topic string, handler func(message interface{}) error) error {
	er := make(chan error)
	go func(err chan error) {
		subscribtion, e := b.connection.Subscribe(topic, func(msg *stan.Msg) {
			error := handler(msg.Data)
			if error == nil {
				msg.Ack()
				log.Debug.Println("Send ACK")
			}
		},
			stan.DurableName(clientId),
			stan.StartWithLastReceived(), stan.SetManualAckMode())

		b.subscribtions[clientId] = func() error {
			return subscribtion.Unsubscribe()
		}
		if e != nil {
			err <- e
		}
	}(er)
	return <-er
}
func (b *natsBroker) Unsubscribe(clientId string) error {
	unsubscribe := b.subscribtions[clientId]
	err := unsubscribe()
	delete(b.subscribtions, clientId)

	return err
}
