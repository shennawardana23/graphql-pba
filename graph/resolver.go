package graph

import (
	_ "github.com/lib/pq"
	"github.com/shennawardana23/graphql-pba/internal/service"
)

// Resolver serves as dependency injection container
type Resolver struct {
	UserService *service.UserService
}
