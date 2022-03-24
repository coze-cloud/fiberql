package fibergraphql

import (
	"context"

	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/graphql"
)

type Subscriber struct {
	OperationId string
	operation   graphql.Params
	ctx         context.Context
	cancelFunc  context.CancelFunc
	conn        *websocket.Conn
}

func NewSubscriber(
	ctx context.Context,
	schema graphql.Schema,
	conn *websocket.Conn,
	operationId string,
	operation map[string]interface{},
) *Subscriber {
	subscriberCtx, subscriberCtxCancel := context.WithCancel(ctx)

	operationName := ""
	if operation["operationName"] != nil {
		operationName = operation["operationName"].(string)
	}

	graphqlOperation := graphql.Params{
		Context:        ctx,
		Schema:         schema,
		VariableValues: operation["variables"].(map[string]interface{}),
		OperationName:  operationName,
		RequestString:  operation["query"].(string),
	}

	subscriber := &Subscriber{
		ctx:         subscriberCtx,
		cancelFunc:  subscriberCtxCancel,
		conn:        conn,
		OperationId: operationId,
		operation:   graphqlOperation,
	}

	go subscriber.execute()

	return subscriber
}

func (s *Subscriber) execute() {
	subscribeChannel := graphql.Subscribe(s.operation)

	for {
		select {
		case <-s.ctx.Done():
			return
		case result, ok := <-subscribeChannel:
			if !ok {
				s.conn.WriteJSON(NewCompleteMessage(
					s.OperationId,
				))
				return
			}
			if err := s.sendResult(result); err != nil {
				return
			}
		}
	}
}

func (s *Subscriber) sendResult(result *graphql.Result) error {
	if result.HasErrors() {
		return s.conn.WriteJSON(NewErrorMessage(
			s.OperationId,
			result.Errors,
		))
	}

	return s.conn.WriteJSON(NewNextMessage(
		s.OperationId,
		result.Data,
	))
}

func (s *Subscriber) Unsubscribe() {
	s.cancelFunc()
}
