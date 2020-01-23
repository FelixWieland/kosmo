package kosmo

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type metaInformations struct {
	name        string
	description string
}

type functionInformations struct {
	metaInformations
	args     graphql.FieldConfigArgument
	resolver func(graphql.ResolveParams) (interface{}, error)
}

type typeInformations struct {
	metaInformations
	typ graphql.Output
}

// INFORMATIONS
// Contains all information needed the build a graphql schema through reflection
// metaInformations name and description are not always both filled

func reflectFunctionInformations(function interface{}) functionInformations {
	function, description := validateResolverArgument(function)
	reflectedFunction := reflect.ValueOf(function)

	name := runtimeFunctionName(reflectedFunction)
	args := functionToConfigArguments(reflectedFunction)
	resolver := functionToResolver(reflectedFunction)

	return functionInformations{
		metaInformations: metaInformations{
			name:        name,
			description: description,
		},
		args:     args,
		resolver: resolver,
	}
}

func reflectTypeInformations(value interface{}) typeInformations {
	value, description := validateResolverArgument(value)
	rinfos := typeInformations{
		metaInformations: metaInformations{
			description: description,
		},
	}

	switch reflectTypeKind(value) {
	case "struct":
		rinfos.typ = structToGraph(value)
	default:
		rinfos.typ = sliceToGraph(value)
	}

	return rinfos
}

// RESOLVER
// uses reflection to extract the arguments of a given function and creates a resolver type based
// on that arguments. The functionToResolver func returns a resolver that can be used in graphql-go's
// graphql library. The resolver-function call is currently kind of expensive becaus it uses reflection
// to call the passed in function

func functionToConfigArguments(fn reflect.Value) graphql.FieldConfigArgument {
	argumentConfig := graphql.FieldConfigArgument{}
	arg, hasArgs := reflectArgumentFromResolverFunction(fn)
	if !hasArgs {
		return argumentConfig
	}

	for i := 0; i < arg.NumField(); i++ {
		argField := arg.Field(i)
		argumentConfig[argField.Name] = &graphql.ArgumentConfig{
			Type: nativeTypeToGraphQL(argField.Type.Name()),
		}
	}

	return argumentConfig
}

func functionToResolver(fn reflect.Value) func(graphql.ResolveParams) (interface{}, error) {
	arg, _ := reflectArgumentFromResolverFunction(fn)
	return func(p graphql.ResolveParams) (interface{}, error) {

		functionArguments := []reflect.Value{}
		if arg != nil {
			functionArguments = []reflect.Value{createFunctionStructArgumentFromMap(arg, p.Args)}
		}

		results := fn.Call(functionArguments)

		returnValue := results[0].Interface()
		returnError := results[1].Interface()

		if returnError != nil {
			return returnValue, returnError.(error)
		}
		return returnValue, nil
	}
}

// TYPES
// uses reflection to extract type informations of the given value.
// the created graphql.Objects are cached to prevent the "Multiple types with the same name" error
// that graphql-go's graphql library returns if multiple graphql.Objects with the same name are created

func structToGraph(genStruct interface{}) *graphql.Object {
	conf := structToGraphConfig(genStruct)
	obj := graphqlObjectCache.Read(conf.Name, gqlObjFallbackFactory(conf)).(*graphql.Object)
	return obj
}

func sliceToGraph(genSlice interface{}) *graphql.List {
	conf := sliceToGraphConfig(genSlice)
	obj := graphqlObjectCache.Read(conf.Name, gqlObjFallbackFactory(conf)).(*graphql.Object)
	return graphql.NewList(obj)
}

// CONFIGS
// the created graphql.ObjectConfigs are build from the field names of the underling struct type
// the created configs are used to create graphql.Objects

func structToGraphConfig(genStruct interface{}) graphql.ObjectConfig {
	underlingType := reflect.TypeOf(genStruct)
	return buildObjectConfigFromType(underlingType)
}

func sliceToGraphConfig(genSlice interface{}) graphql.ObjectConfig {
	underlingType := reflect.TypeOf(genSlice).Elem()
	return buildObjectConfigFromType(underlingType)
}

func buildObjectConfigFromType(reflectedType reflect.Type) graphql.ObjectConfig {
	name, infos := reflectStructInformations(reflectedType)
	fields := nativeFieldsToGraphQLFields(infos)
	return graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	}
}

// FIELDS
// Uses reflection to map the fields in a struct-type to create graphql.Types
// based on the type of the field

func nativeFieldsToGraphQLFields(fields []reflect.StructField) graphql.Fields {
	graphQLFields := graphql.Fields{}
	for _, field := range fields {
		typ := nativeFieldToGraphQL(field)
		graphQLFields[field.Name] = &typ
	}
	return graphQLFields
}

func nativeFieldToGraphQL(field reflect.StructField) graphql.Field {
	var nTyp graphql.Output

	b := field.Type.Kind().String()
	_ = b

	switch field.Type.Kind().String() {
	case "struct":
		conf := buildObjectConfigFromType(field.Type)
		nTyp = graphqlObjectCache.Read(conf.Name, gqlObjFallbackFactory(conf)).(*graphql.Object)
	case "slice":
		conf := buildObjectConfigFromType(field.Type.Elem())
		nTyp = graphql.NewList(graphqlObjectCache.Read(conf.Name, gqlObjFallbackFactory(conf)).(*graphql.Object))
	default:
		nTyp = nativeTypeToGraphQL(field.Type.Name())
	}

	return graphql.Field{
		Type:        nTyp,
		Description: field.Tag.Get("description"),
	}
}

func nativeTypeToGraphQL(typeName string) graphql.Type {
	switch typeName {
	case "int":
		return graphql.Int
	case "uint":
		return graphql.Int
	case "string":
		return graphql.String
	case "float32":
		return graphql.Float
	case "float64":
		return graphql.Float
	default:
		panic(typeName + " is not implemented yet")
	}
}
