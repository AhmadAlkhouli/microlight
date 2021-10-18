package message

type Command struct {
	//ID            interface{}      `json:"id"`
	CommandType   string      `json:"command_type"`
	AggregateID   interface{} `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	CommandData   interface{} `json:"command_data"`
}

/* //Command ...
type Command interface {
	//GetID() interface{}
	GetCommandType() string
	GetAggregateID() interface{}
	GetAggregateType() string
	GetCommandData() interface{}
} */

//CreateCommand ...
func CreateCommand(
	CommandType string,
	AggregateID interface{},
	AggregateType string,
	CommandData interface{}) *Command {
	com := Command{}
	//com.ID = ID
	com.CommandType = CommandType
	com.AggregateID = AggregateID
	com.AggregateType = AggregateType
	com.CommandData = CommandData
	return &com
}
