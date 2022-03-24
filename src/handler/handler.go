package handler

import (
	"github.com/coze-cloud/fiberql/src/handler/graphiql"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type handler struct {
	Schema   graphql.Schema
	GraphiQl bool
}

func NewHandler(config Config) *handler {
	return &handler{
		Schema:   config.Schema,
		GraphiQl: config.GraphiQl,
	}
}

func (h *handler) Handle(ctx *fiber.Ctx) error {
	raw := ctx.Context().QueryArgs().Has("raw")
	method := ctx.Method()
	if h.GraphiQl && method == fiber.MethodGet && !raw {
		return graphiql.NewHandler().Handle(ctx)
	}

	request, err := newRequest(ctx)
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
