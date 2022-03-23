package main

import (
	"time"

	fibergraphql "github.com/coze-hosting/fiber-graphql/src"
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

	subscriptionHandler := fibergraphql.NewSubscriptionHandler(
		schema,
		[]string{"*"},
		time.Second*10,
		time.Minute,
	)
	app.Get("/subscriptions", subscriptionHandler.Handle)

	handler := fibergraphql.NewHandler(
		schema,
	)
	app.Get("/graphql", handler.Handle)
	app.Post("/graphql", handler.Handle)

	app.Listen(":3000")
}
