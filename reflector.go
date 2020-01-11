package kosmo

import (
	"reflect"

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

func nativeFieldToGraphQL(field reflect.StructField) graphql.Field {
	return graphql.Field{
		Type: func() graphql.Type {
			switch field.Type.Name() {
			case "int":
				return graphql.Int
			case "string":
				return graphql.String
			case "float":
				return graphql.Float
			default:
				return graphql.String
			}
		}(),
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
