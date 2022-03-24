package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	Query         string                 `query:"query" json:"query" form:"query"`
	Variables     map[string]interface{} `query:"variables" json:"variables" form:"variables"`
	OperationName string                 `query:"operationName" json:"operationName" form:"operationName"`
}

func newRequest(ctx *fiber.Ctx) (*request, *fiber.Error) {
	if ctx.Query("query") != "" {
		return newRequestFromQuery(ctx)
	}
	return newRequestFromBody(ctx)
}

func newRequestFromQuery(ctx *fiber.Ctx) (*request, *fiber.Error) {
	request := &request{}
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

func newRequestFromBody(ctx *fiber.Ctx) (*request, *fiber.Error) {
	body := ctx.Body()

	contentType := string(ctx.Request().Header.ContentType())
	switch contentType {
	case "application/graphql":
		return &request{
			Query: string(body),
		}, nil
	case "application/json":
		request := &request{}
		if err := ctx.BodyParser(request); err != nil {
			return nil, fiber.ErrBadRequest
		}
		return request, nil
	case "application/x-www-form-urlencoded":
		request := &request{}
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
