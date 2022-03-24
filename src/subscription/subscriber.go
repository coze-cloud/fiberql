package subscription

import (
	"context"

	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/graphql"
)

type subscriber struct {
	OperationId string
	operation   graphql.Params
	ctx         context.Context
	cancelFunc  context.CancelFunc
	conn        *websocket.Conn
}

func newSubscriber(
	ctx context.Context,
	schema graphql.Schema,
	conn *websocket.Conn,
	operationId string,
	operation map[string]interface{},
) *subscriber {
	subscriberCtx, subscriberCtxCancel := context.WithCancel(ctx)

	operationName := ""
	if operation["operationName"] != nil {
		operationName = operation["operationName"].(string)
	}

	variables := map[string]interface{}{}
	if operation["variables"] != nil {
		variables = operation["variables"].(map[string]interface{})
	}

	graphqlOperation := graphql.Params{
		Context:        ctx,
		Schema:         schema,
		VariableValues: variables,
		OperationName:  operationName,
		RequestString:  operation["query"].(string),
	}

	subscriber := &subscriber{
		ctx:         subscriberCtx,
		cancelFunc:  subscriberCtxCancel,
		conn:        conn,
		OperationId: operationId,
		operation:   graphqlOperation,
	}

	go subscriber.execute()

	return subscriber
}

func (s *subscriber) execute() {
	subscribeChannel := graphql.Subscribe(s.operation)

	for {
		select {
		case <-s.ctx.Done():
			return
		case result, ok := <-subscribeChannel:
			if !ok {
				s.conn.WriteJSON(newCompleteMessage(
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

func (s *subscriber) sendResult(result *graphql.Result) error {
	if result.HasErrors() {
		return s.conn.WriteJSON(newErrorMessage(
			s.OperationId,
			result.Errors,
		))
	}

	return s.conn.WriteJSON(newNextMessage(
		s.OperationId,
		result.Data,
	))
}

func (s *subscriber) Unsubscribe() {
	s.cancelFunc()
}
