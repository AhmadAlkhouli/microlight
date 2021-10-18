package message

// Event ...
type Event struct {
	ID            string      `json:"id"`
	EventType     string      `json:"event_type"`
	AggregateID   interface{} `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	EventData     interface{} `json:"event_data"`
}
