package kosmo

import "github.com/graphql-go/graphql"

// Service represents a kosmo-microservice
type Service struct {
	HTTPConfig HTTPConfig
	graphQL    graphQL
}

// HTTPConfig represents the http configurations
type HTTPConfig struct {
	Port int
}

// ResolveParams Params for Field.resolve()
type ResolveParams graphql.ResolveParams

// Field resolve the given type
type Field interface {
	resolve(struct{}) (interface{}, error)
}

type graphQL graphql.SchemaConfig

// Queries - Adds the queries to the grapql schema
func (s *Service) Queries(types ...interface{}) *Service {
	// s.graphQL.Query :=
	return s
}

// Mutations - Adds the mutations to the graphql schema
func (s *Service) Mutations(types ...interface{}) *Service {
	return s
}

// Start - Starts the http server
func (s *Service) Start() {

}
