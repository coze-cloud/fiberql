package subscription

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/graphql-go/graphql"
)

type handler struct {
	Schema                    graphql.Schema
	Origins                   []string
	ConnectionInitWaitTimeout time.Duration
	PingInterval              time.Duration
}

func NewHandler(config Config) *handler {
	return &handler{
		Schema:                    config.Schema,
		Origins:                   config.Origins,
		ConnectionInitWaitTimeout: config.ConnectionInitWaitTimeout,
		PingInterval:              config.PingInterval,
	}
}

func (h *handler) Handle(ctx *fiber.Ctx) error {
	return websocket.New(h.handleWebsocket, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Origins:         h.Origins,
		Subprotocols:    []string{"graphql-transport-ws"},
	})(ctx)
}

func (h *handler) handleWebsocket(conn *websocket.Conn) {
	subconn := newConnection(conn, h.Schema)
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
