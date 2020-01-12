package kosmo

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type metaInformations struct {
	name        string
	description string
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
	underlingType := reflect.TypeOf(genSlice).Elem()
	fields := fieldsFromType(underlingType)
	object := graphql.NewObject(graphql.ObjectConfig{
		Name:   underlingType.Name(),
		Fields: fields,
	})
	return graphql.NewList(object)
}

func buildQueryField(object *graphql.Object, args graphql.FieldConfigArgument, resolver func(graphql.ResolveParams) (interface{}, error)) graphql.Field {
	return graphql.Field{
		Type:    object,
		Args:    args,
		Resolve: resolver,
	}
}

func reflectResolverMethod(genStruct interface{}) reflect.Value {
	t := reflect.ValueOf(genStruct)
	method := t.MethodByName("Resolve")
	empty := reflect.Value{}

	if method == empty {
		panic(reflect.TypeOf(t).Name() + " does not implement resolve(...interface{}) (interface{}, error)")
	}

	return method
}

func getArgumentFromResolverMethod(method reflect.Value) (reflect.Type, bool) {
	if method.Type().NumIn() == 0 {
		return nil, false
	}
	return method.Type().In(0), true
}

func reflectArgsFromResolver(method reflect.Value) graphql.FieldConfigArgument {
	argumentConfig := graphql.FieldConfigArgument{}

	arg, hasArgs := getArgumentFromResolverMethod(method)
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

func resolverFactory(method reflect.Value) func(graphql.ResolveParams) (interface{}, error) {
	arg, _ := getArgumentFromResolverMethod(method)

	return func(p graphql.ResolveParams) (interface{}, error) {
		raw := reflect.New(arg).Elem()

		for key, field := range p.Args {
			raw.FieldByName(key).Set(reflect.ValueOf(field))
		}

		functionArguments := []reflect.Value{raw}

		results := method.Call(functionArguments)
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
