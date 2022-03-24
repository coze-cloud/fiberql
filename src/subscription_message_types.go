package fibergraphql

import "github.com/graphql-go/graphql/gqlerrors"

const (
	NormalClosure                   = 1000
	ConnectionInitialisationTimeout = 4408
	TooManyInitialisationRequests   = 4429
	SubscriberAlreadyExists         = 4409
	Unauthorized                    = 4401
	InvalidMessage                  = 4400
)

const (
	ConnectionInit = "connection_init"
	ConnectionAck  = "connection_ack"
	Ping           = "ping"
	Pong           = "pong"
	Subscribe      = "subscribe"
	Next           = "next"
	Error          = "error"
	Complete       = "complete"
)

type ConnectionMessage struct {
	Type    string                 `json:"type"`
	Id      string                 `json:"id,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func NewConnectionAckMessage() *ConnectionMessage {
	return &ConnectionMessage{
		Type: ConnectionAck,
	}
}

func NewPingMessage() *ConnectionMessage {
	return &ConnectionMessage{
		Type: Ping,
	}
}

func NewPongMessage() *ConnectionMessage {
	return &ConnectionMessage{
		Type: Pong,
	}
}

func NewNextMessage(id string, result interface{}) *ConnectionMessage {
	return &ConnectionMessage{
		Type: Next,
		Id:   id,
		Payload: map[string]interface{}{
			"data": result,
		},
	}
}

func NewErrorMessage(id string, errors []gqlerrors.FormattedError) *ConnectionMessage {
	return &ConnectionMessage{
		Type: Error,
		Id:   id,
		Payload: map[string]interface{}{
			"errors": errors,
		},
	}
}

func NewCompleteMessage(id string) *ConnectionMessage {
	return &ConnectionMessage{
		Type: Complete,
		Id:   id,
	}
}
