package graphiql

import (
	"github.com/gofiber/fiber/v2"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (*handler) Handle(ctx *fiber.Ctx) error {
	protocol := "ws://"
	if ctx.Protocol() == "https" {
		protocol = "wss://"
	}
	url := ctx.OriginalURL()
	subscriptionUrl := protocol + string(ctx.Context().Host()) + "/subscriptions"

	template := graphiQlTemplate(
		url,
		subscriptionUrl,
	)

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return ctx.Status(fiber.StatusOK).SendString(template)
}
