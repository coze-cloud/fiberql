package fibergraphql

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type HandlerRequest struct {
	Query         string                 `query:"query" json:"query" form:"query"`
	Variables     map[string]interface{} `query:"variables" json:"variables" form:"variables"`
	OperationName string                 `query:"operationName" json:"operationName" form:"operationName"`
}

func NewHandlerRequest(ctx *fiber.Ctx) (*HandlerRequest, *fiber.Error) {
	if ctx.Query("query") != "" {
		return newHandlerRequestFromQuery(ctx)
	}
	return newHandlerRequestFromBody(ctx)
}

func newHandlerRequestFromQuery(ctx *fiber.Ctx) (*HandlerRequest, *fiber.Error) {
	request := &HandlerRequest{}
	if err := ctx.QueryParser(request); err != nil {
		return nil, fiber.ErrBadRequest
	}
	variables := ctx.Query("variables")
	if len(variables) <= 0 {
		return request, nil
	}
	if err := json.Unmarshal([]byte(variables), &request.Variables); err != nil {
		return nil, fiber.ErrBadRequest
	}
	return request, nil
}

func newHandlerRequestFromBody(ctx *fiber.Ctx) (*HandlerRequest, *fiber.Error) {
	body := ctx.Body()

	contentType := string(ctx.Request().Header.ContentType())
	switch contentType {
	case "application/graphql":
		return &HandlerRequest{
			Query: string(body),
		}, nil
	case "application/json":
		request := &HandlerRequest{}
		if err := ctx.BodyParser(request); err != nil {
			return nil, fiber.ErrBadRequest
		}
		return request, nil
	case "application/x-www-form-urlencoded":
		request := &HandlerRequest{}
		if err := ctx.BodyParser(request); err != nil {
			return nil, fiber.ErrBadRequest
		}
		variables := ctx.FormValue("variables")
		if len(variables) <= 0 {
			return request, nil
		}
		if err := json.Unmarshal([]byte(variables), &request.Variables); err != nil {
			return nil, fiber.ErrBadRequest
		}
		return request, nil
	default:
		return nil, fiber.ErrBadRequest
	}
}
