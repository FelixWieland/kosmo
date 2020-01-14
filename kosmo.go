package kosmo

import (
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
	Port    string
	APIBase string
}

// GraphQLConfig represents the graphql configuraiton options
type GraphQLConfig struct {
	UseTypeAsQueryName      bool
	ReplaceResolverPrefixes bool
	ResolverPrefixes        []string
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

// Type creates a graphQL Schema
func Type(typedVar interface{}) GraphQLSchema {
	if reflectIsFromTypeSlice(typedVar) {
		return GraphQLSchema{
			resolverType: sliceToGraph(typedVar),
		}
	}

	return GraphQLSchema{
		resolverType: structToGraph(typedVar),
	}
}

// Query adds the Query resolver to the Type
func (t GraphQLSchema) Query(resolverFunction interface{}) GraphQLSchema {
	fnInfos := reflectFunctionInformations(resolverFunction)

	t.query = graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			fnInfos.name: &graphql.Field{
				Type:    t.resolverType,
				Args:    fnInfos.args,
				Resolve: fnInfos.resolver,
			},
		},
	}

	return t
}

// Mutations adds the GraphQLMutation functions
func (t *GraphQLSchema) Mutations(resolverFunctions ...interface{}) *GraphQLSchema {
	fields := graphql.Fields{}

	for key := range resolverFunctions {
		fnInfos := reflectFunctionInformations(resolverFunctions[key])

		fields[fnInfos.name] = &graphql.Field{
			Type:    t.resolverType,
			Args:    fnInfos.args,
			Resolve: fnInfos.resolver,
		}
	}

	t.mutations = graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: fields,
	}
	return t
}

// Schemas adds the schemas to the service
func (s *Service) Schemas(schemas ...GraphQLSchema) *Service {
	queryConfigs := []graphql.ObjectConfig{}
	mutationConfigs := []graphql.ObjectConfig{}

	for _, schema := range schemas {

		if s.GraphQLConfig.UseTypeAsQueryName {
			fields := schema.query.Fields.(graphql.Fields)
			schema.query.Fields = rewriteFieldNamesToType(fields)
		}

		if s.GraphQLConfig.ReplaceResolverPrefixes {
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

// Start - Starts the http server
func (s *Service) Start() error {
	schema, err := graphql.NewSchema(s.graphQL)

	if err != nil {
		panic(err)
	}

	return startServer(s.HTTPConfig, schema)
}
