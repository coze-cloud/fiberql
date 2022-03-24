package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/graphql"
)

type connection struct {
	conn        *websocket.Conn
	schema      graphql.Schema
	initialized bool
	lastPong    time.Time
	subscribers map[string]*subscriber
}

func newConnection(conn *websocket.Conn, schema graphql.Schema) *connection {
	return &connection{
		conn:        conn,
		schema:      schema,
		initialized: false,
		lastPong:    time.Time{},
		subscribers: make(map[string]*subscriber),
	}
}

func (c *connection) ConnectionInitialisationTimeout(ctx context.Context, timeout time.Duration) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeoutCtx.Done():
			if c.initialized {
				return
			}
			c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
				ConnectionInitialisationTimeout,
				"Connection initialisation timeout",
			))
			return
		}
	}
}

func (c *connection) Ping(ctx context.Context, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			if !c.initialized {
				continue
			}
			c.conn.WriteJSON(newPingMessage())
		}
	}
}

func (c *connection) Handle(ctx context.Context, message *ConnectionMessage) {
	c.handleConnectionInit(message)
	c.handlePing(message)
	c.handlePong(message)

	c.handleSubscribe(ctx, message)
	c.handleComplete(message)
}

func (c *connection) handleConnectionInit(message *ConnectionMessage) {
	if message.Type != ConnectionInit {
		return
	}

	if c.initialized {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			TooManyInitialisationRequests,
			"Too many initialisation requests",
		))
		return
	}

	c.conn.WriteJSON(newConnectionAckMessage())

	c.initialized = true
}

func (c *connection) handlePing(message *ConnectionMessage) {
	if message.Type != Ping {
		return
	}

	c.conn.WriteJSON(newPongMessage())
}

func (c *connection) handlePong(message *ConnectionMessage) {
	if message.Type != Pong {
		return
	}

	c.lastPong = time.Now().UTC()
}

func (c *connection) handleSubscribe(ctx context.Context, message *ConnectionMessage) {
	if message.Type != Subscribe {
		return
	}

	if !c.initialized {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			Unauthorized,
			"Unauthorized",
		))
		return
	}

	if subscriber, ok := c.subscribers[message.Id]; ok {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			SubscriberAlreadyExists,
			fmt.Sprintf("Subscriber for %s already exists", subscriber.OperationId),
		))
		return
	}

	subscriber := newSubscriber(
		ctx,
		c.schema,
		c.conn,
		message.Id,
		message.Payload,
	)
	c.subscribers[message.Id] = subscriber
}

func (c *connection) handleComplete(message *ConnectionMessage) {
	if message.Type != Complete {
		return
	}

	if !c.initialized {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
			Unauthorized,
			"Unauthorized",
		))
		return
	}

	subscriber, ok := c.subscribers[message.Id]
	if !ok {
		return
	}

	subscriber.Unsubscribe()

	delete(c.subscribers, message.Id)
}
