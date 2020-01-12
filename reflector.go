package kosmo

import (
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/graphql-go/graphql"
)

func reflectGormType(gormschema interface{}) interface{} {
	_, fields := describeStruct(gormschema)
	return fields
}

func describeStruct(genStruct interface{}) (string, []reflect.StructField) {
	t := reflect.TypeOf(genStruct)
	fields := []reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i))
	}
	return t.Name(), fields
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

func structToGraph(genStruct interface{}) *graphql.Object {
	name, infos := describeStruct(genStruct)
	fields := graphql.Fields{}
	for _, field := range infos {
		typ := nativeFieldToGraphQL(field)
		fields[field.Name] = &typ
	}
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   name,
			Fields: fields,
		},
	)
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
	_ = arg
	return func(p graphql.ResolveParams) (interface{}, error) {
		raw := reflect.New(arg).Elem().Interface()
		err := mapstructure.Decode(p.Args, &raw)
		if err != nil {
			panic(err)
		}

		inputs := []reflect.Value{
			reflect.ValueOf(raw),
		}

		returnValues := method.Call(inputs)
		returnValue := returnValues[0].Interface()
		returnError := returnValues[1].Interface()

		if returnError != nil {
			return returnValue, returnError.(error)
		}
		return returnValue, nil
	}
}
