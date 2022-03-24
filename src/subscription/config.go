package subscription

import (
	"time"

	"github.com/graphql-go/graphql"
)

type Config struct {
	Schema                    graphql.Schema
	Origins                   []string
	ConnectionInitWaitTimeout time.Duration
	PingInterval              time.Duration
}
