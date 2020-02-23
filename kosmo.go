package kosmo

import (
	"net/http"

	"github.com/graphql-go/graphql"
)

// Service represents a kosmo-microservice
type Service struct {
	HTTPConfig    HTTPConfig
	GraphQLConfig GraphQLConfig
	graphQL       graphql.SchemaConfig
}

// HTTPConfig represents the http configurations
type HTTPConfig struct {
	Port       string
	APIBase    string
	Playground bool
	Gzip       bool
}

// GraphQLConfig represents the graphql configuration options
type GraphQLConfig struct {
	RemoveResolverPrefixes bool
	ResolverPrefixes       []string
}

// ResolveParams Params for Field.resolve()
type ResolveParams graphql.ResolveParams

// Field resolve the given type
type Field interface {
	resolve(struct{}) (interface{}, error)
}

// GraphQLSchema contains the type, the query and the mutations of the type
type GraphQLSchema struct {
	resolverType graphql.Output
	query        graphql.ObjectConfig
	mutations    graphql.ObjectConfig
}

// Describer is a wrapper for any type that allows a description
type Describer struct {
	Value       interface{}
	Description string
}

// Type creates a graphQL Schema
func Type(typedVar interface{}) *GraphQLSchema {
	typeInfos := reflectTypeInformations(typedVar)

	return &GraphQLSchema{
		resolverType: typeInfos.typ,
	}
}

// Queries adds the Query resolver to the Type
func (t *GraphQLSchema) Queries(resolverFunctions ...interface{}) *GraphQLSchema {
	fields := graphql.Fields{}

	for key := range resolverFunctions {
		functionInfos := reflectFunctionInformations(resolverFunctions[key])

		fields[functionInfos.name] = &graphql.Field{
			Type:        t.resolverType,
			Args:        functionInfos.args,
			Resolve:     functionInfos.resolver,
			Description: functionInfos.description,
		}
	}

	t.query = graphql.ObjectConfig{
		Name:        "Query",
		Description: "Root for all queries",
		Fields:      fields,
	}
	return t
}

// Mutations adds the GraphQLMutation functions
func (t *GraphQLSchema) Mutations(resolverFunctions ...interface{}) *GraphQLSchema {
	fields := graphql.Fields{}

	for key := range resolverFunctions {
		functionInfos := reflectFunctionInformations(resolverFunctions[key])

		fields[functionInfos.name] = &graphql.Field{
			Type:        t.resolverType,
			Args:        functionInfos.args,
			Resolve:     functionInfos.resolver,
			Description: functionInfos.description,
		}
	}

	t.mutations = graphql.ObjectConfig{
		Name:        "Mutation",
		Description: "Root for all mutations",
		Fields:      fields,
	}
	return t
}

// Schemas adds the schemas to the service
func (s *Service) Schemas(schemas ...*GraphQLSchema) *Service {
	queryConfigs := []graphql.ObjectConfig{}
	mutationConfigs := []graphql.ObjectConfig{}

	for _, schema := range schemas {

		if s.GraphQLConfig.RemoveResolverPrefixes {
			fields := schema.query.Fields.(graphql.Fields)
			schema.query.Fields = replaceResolverPrefixes(s.GraphQLConfig.ResolverPrefixes, fields)
		}

		queryConfigs = append(queryConfigs, schema.query)
		mutationConfigs = append(mutationConfigs, schema.mutations)
	}

	s.graphQL = graphql.SchemaConfig{
		Query:    makeGraphQLObject(combineObjectConfig(queryConfigs...)),
		Mutation: makeGraphQLObject(combineObjectConfig(mutationConfigs...)),
	}

	return s
}

// Server - returns the http server
func (s *Service) Server() *http.Server {
	schema, err := graphql.NewSchema(s.graphQL)
	if err != nil {
		panic(err)
	}

	return muxServer(s.HTTPConfig, schema)
}

//Handler - returns the raw HTTP handler
func (s *Service) Handler() http.Handler {
	schema, err := graphql.NewSchema(s.graphQL)
	if err != nil {
		panic(err)
	}

	return rawHandler(s.HTTPConfig.Playground, schema)
}
