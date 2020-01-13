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

func reflectGormType(gormschema interface{}) interface{} {
	_, fields := describeStruct(gormschema)
	return fields
}

func fieldsFromType(t reflect.Type) []reflect.StructField {
	fields := []reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i))
	}
	return fields
}

func describeStruct(genStruct interface{}) (string, []reflect.StructField) {
	t := reflect.TypeOf(genStruct)
	return t.Name(), fieldsFromType(t)
}

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

func nativeFieldToGraphQL(field reflect.StructField) graphql.Field {
	return graphql.Field{
		Type: nativeTypeToGraphQL(field.Type.Name()),
	}
}

func reflectedFieldsToGraphQL(fields []reflect.StructField) *graphql.Fields {
	graphQLFields := graphql.Fields{}
	for _, field := range fields {
		typ := nativeFieldToGraphQL(field)
		graphQLFields[field.Name] = &typ
	}
	return &graphQLFields
}

func structToGraph(genStruct interface{}) *graphql.Object {
	name, infos := describeStruct(genStruct)
	fields := reflectedFieldsToGraphQL(infos)
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		},
	)
}

func sliceToGraph(genSlice interface{}) *graphql.List {
	return graphql.NewList(reflectGraphTypeFromSlice(genSlice))
}

func reflectGraphTypeFromSlice(genSlice interface{}) *graphql.Object {
	underlingType := reflect.TypeOf(genSlice).Elem()
	fields := fieldsFromType(underlingType)
	object := graphql.NewObject(graphql.ObjectConfig{
		Name:   underlingType.Name(),
		Fields: reflectedFieldsToGraphQL(fields),
	})
	return object
}

func buildQueryField(object *graphql.Object, args graphql.FieldConfigArgument, resolver func(graphql.ResolveParams) (interface{}, error)) graphql.Field {
	return graphql.Field{
		Type:    object,
		Args:    args,
		Resolve: resolver,
	}
}

func reflectResolverFunction(fn interface{}) reflect.Value {
	t := reflect.ValueOf(fn)
	return t
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

func reflectMetaInformations(genVar interface{}) metaInformations {
	return metaInformations{
		name: reflect.TypeOf(genVar).Name(),
	}
}

func reflectFunctionInformations(fn interface{}) functionInformations {
	reflectedFunction := reflectResolverFunction(fn)

	name := getFunctionName(fn)
	args := reflectArgsFromResolver(reflectedFunction)
	resolver := resolverFactory(reflectedFunction)

	return functionInformations{
		metaInformations: metaInformations{
			name: name,
		},
		args:     args,
		resolver: resolver,
	}
}

func reflectIsFromTypeSlice(genVar interface{}) bool {
	return reflect.TypeOf(genVar).Kind().String() != "struct"
}

func getFunctionName(i interface{}) string {
	nameWithPackage := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.SplitAfter(nameWithPackage, ".")

	return parts[len(parts)-1]
}
