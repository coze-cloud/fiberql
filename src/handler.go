package fibergraphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type Handler struct {
	Schema graphql.Schema
}

func NewHandler(schema graphql.Schema) *Handler {
	return &Handler{
		Schema: schema,
	}
}

func (h *Handler) Handle(ctx *fiber.Ctx) error {
	request, err := NewHandlerRequest(ctx)
	if err != nil {
		return err
	}

	params := graphql.Params{
		Schema:         h.Schema,
		RequestString:  request.Query,
		VariableValues: request.Variables,
		OperationName:  request.OperationName,
	}

	result := graphql.Do(params)

	return ctx.Status(fiber.StatusOK).JSON(result)
}
