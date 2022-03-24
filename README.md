# fiberql
🔦 Simple GraphQL handler for fiber

## Installation

Adding fiberql to your Go module is as easy as calling this command in your project 

```shell
go get github.com/coze-cloud/fiberql
```

**Important:** For fiberql you have to use [gofiber/fiber](https://github.com/gofiber/fiber) as webserver and [graphql-go/graphql](https://github.com/graphql-go/graphql) for your schema.

## Hooking it up

```go
// Have a look at the graphql-go/graphql docs for further
// information about defining the schema 
schema, _ := graphql.NewSchema(...)

app := fiber.New()

graphql := handler.NewHandler(handler.Config{
    Schema:   schema,
    GraphiQl: true,
})
app.Get("/graphql", graphql.Handle)
app.Post("/graphql", graphql.Handle)

// The subscription handler serves a websocket 
// communicating via the graphql-transport-ws protocol
subscriptions := subscription.NewHandler(subscription.Config{
    Schema:                    schema,
    Origins:                   []string{"*"},
    ConnectionInitWaitTimeout: time.Second,
    PingInterval:              time.Minute,
})
app.Get("/subscriptions", subscriptions.Handle)

app.Listen(":3000")
```


---

Copyright © 2022 - The cozy team **& contributors**
