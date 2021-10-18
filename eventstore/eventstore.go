package eventstore

import (
	"gopkg.in/mgo.v2"
)

const (
	//URI ...
	URI = "mongodb://localhost:27017"
)

type eventStore struct {
	URL        string
	DBName     string
	Collection string
}
type EventStore interface {
	Persist(event interface{})
}

//NewEventStore ...
func NewEventStore(url string, dbname string, collection string) EventStore {
	return &eventStore{URL: url, DBName: dbname, Collection: collection}
}

//Persist ...
func (e *eventStore) Persist(event interface{}) {

	session, err := mgo.Dial(e.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	col := session.DB(e.DBName).C(e.Collection)
	col.Insert(event)
}
