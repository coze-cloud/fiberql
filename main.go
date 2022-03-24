package main

import (
	"time"

	"github.com/coze-cloud/fiberql/src/handler"
	"github.com/coze-cloud/fiberql/src/subscription"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

func main() {
	fields := graphql.Fields{
		"counter": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source, nil
			},
			Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
				c := make(chan interface{})
				go func() {
					var i int
					for {
						i++
						select {
						case <-p.Context.Done():
							close(c)
							return
						default:
							time.Sleep(time.Second)
							c <- i
						}
						if i == 10 {
							close(c)
							return
						}
					}
				}()
				return c, nil
			},
		},
	}

	rootSubscription := graphql.ObjectConfig{Name: "RootSubscription", Fields: fields}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query:        graphql.NewObject(rootQuery),
		Subscription: graphql.NewObject(rootSubscription),
	}
	schema, _ := graphql.NewSchema(schemaConfig)

	app := fiber.New()

	subHandler := subscription.NewHandler(subscription.Config{
		Schema:                    schema,
		Origins:                   []string{"*"},
		ConnectionInitWaitTimeout: time.Second,
		PingInterval:              time.Minute,
	})
	app.Get("/subscriptions", subHandler.Handle)

	handler := handler.NewHandler(handler.Config{
		Schema:   schema,
		GraphiQl: true,
	})
	app.Get("/graphql", handler.Handle)
	app.Post("/graphql", handler.Handle)

	app.Listen(":3000")
}
