package kosmo

import (
	"reflect"
	"runtime"
	"strings"

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

func fieldsFromType(t reflect.Type) []reflect.StructField {
	fields := []reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i))
	}
	return fields
}

func describeStructType(structType reflect.Type) (string, []reflect.StructField) {
	return structType.Name(), fieldsFromType(structType)
}

// FIELDS

func reflectedFieldsToGraphQL(fields []reflect.StructField) graphql.Fields {
	graphQLFields := graphql.Fields{}
	for _, field := range fields {
		typ := nativeFieldToGraphQL(field)
		graphQLFields[field.Name] = &typ
	}
	return graphQLFields
}

func nativeFieldToGraphQL(field reflect.StructField) graphql.Field {
	return graphql.Field{
		Type:        nativeTypeToGraphQL(field.Type.Name()),
		Description: field.Tag.Get("description"),
	}
}

// maps go's native types to graphql-go's graphql types
func nativeTypeToGraphQL(typeName string) graphql.Type {
	switch typeName {
	case "int":
		return graphql.Int
	case "string":
		return graphql.String
	case "float":
		return graphql.Float
	default:
		return graphql.String
	}
}

// REFLECT CONFIGS -

func structToGraphConfig(genStruct interface{}) graphql.ObjectConfig {
	underlingType := reflect.TypeOf(genStruct)
	return buildObjectConfigFromType(underlingType)
}

func sliceToGraphConfig(genSlice interface{}) graphql.ObjectConfig {
	underlingType := reflect.TypeOf(genSlice).Elem()
	return buildObjectConfigFromType(underlingType)
}

func buildObjectConfigFromType(reflectedType reflect.Type) graphql.ObjectConfig {
	name, infos := describeStructType(reflectedType)
	fields := reflectedFieldsToGraphQL(infos)
	return graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	}
}

// REFLECT TYPES -

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

// REFLECT RESOLVER

func buildQueryField(object *graphql.Object, args graphql.FieldConfigArgument, resolver func(graphql.ResolveParams) (interface{}, error)) graphql.Field {
	return graphql.Field{
		Type:    object,
		Args:    args,
		Resolve: resolver,
	}
}

func getArgumentFromResolverFunction(fn reflect.Value) (reflect.Type, bool) {
	if fn.Type().NumIn() == 0 {
		return nil, false
	}
	return fn.Type().In(0), true
}

func reflectArgsFromResolver(fn reflect.Value) graphql.FieldConfigArgument {
	argumentConfig := graphql.FieldConfigArgument{}

	arg, hasArgs := getArgumentFromResolverFunction(fn)
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

func resolverFactory(fn reflect.Value) func(graphql.ResolveParams) (interface{}, error) {
	arg, _ := getArgumentFromResolverFunction(fn)
	return func(p graphql.ResolveParams) (interface{}, error) {
		raw := reflect.New(arg).Elem()

		for key, field := range p.Args {
			raw.FieldByName(key).Set(reflect.ValueOf(field))
		}

		functionArguments := []reflect.Value{raw}

		results := fn.Call(functionArguments)
		returnValue := results[0].Interface()
		returnError := results[1].Interface()

		if returnError != nil {
			return returnValue, returnError.(error)
		}
		return returnValue, nil
	}
}

func reflectTypeKind(genVar interface{}) string {
	return reflect.TypeOf(genVar).Kind().String()
}

func getFunctionName(function interface{}) string {
	nameWithPackage := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	parts := strings.SplitAfter(nameWithPackage, ".")

	return parts[len(parts)-1]
}

// INFORMATIONS - Endresult

func reflectFunctionInformations(function interface{}) functionInformations {
	function, description := validateResolverArgument(function)
	reflectedFunction := reflect.ValueOf(function)

	name := getFunctionName(function)
	args := reflectArgsFromResolver(reflectedFunction)
	resolver := resolverFactory(reflectedFunction)

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
