package fibergraphql

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/graphql"
)

type SubscriptionHandler struct {
	Schema                    graphql.Schema
	Origins                   []string
	ConnectionInitWaitTimeout time.Duration
	PingInterval              time.Duration
}

func NewSubscriptionHandler(
	schema graphql.Schema,
	origins []string,
	connectionInitWaitTimeout time.Duration,
	pingInterval time.Duration,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		Schema:                    schema,
		Origins:                   origins,
		ConnectionInitWaitTimeout: connectionInitWaitTimeout,
		PingInterval:              pingInterval,
	}
}

func (h *SubscriptionHandler) Handle(ctx *fiber.Ctx) error {
	return websocket.New(h.HandleWebsocket, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Origins:         h.Origins,
		Subprotocols:    []string{"graphql-transport-ws"},
	})(ctx)
}

func (h *SubscriptionHandler) HandleWebsocket(conn *websocket.Conn) {
	subconn := NewSubscriptionConnection(conn, h.Schema)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go subconn.ConnectionInitialisationTimeout(ctx, h.ConnectionInitWaitTimeout)
	go subconn.Ping(ctx, h.PingInterval)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
				InvalidMessage,
				"Invalid message",
			))
			break
		}
		connectionMessage := &ConnectionMessage{}
		if err := json.Unmarshal(message, connectionMessage); err != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
				InvalidMessage,
				"Invalid message",
			))
			break
		}

		subconn.Handle(ctx, connectionMessage)
	}
}
