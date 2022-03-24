package handler

import "github.com/graphql-go/graphql"

type Config struct {
	Schema   graphql.Schema
	GraphiQl bool
}
