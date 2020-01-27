package kosmo

import (
	"testing"

	"github.com/graphql-go/graphql"
	. "github.com/smartystreets/goconvey/convey"
)

type TKosmoStruct struct {
	Name   string
	Number int
}

type TAnotherKosmoStruct struct {
	Name string
}

func TResolveKosmoStruct() (TKosmoStruct, error) {
	return TKosmoStruct{
		Name:   "Test",
		Number: 1,
	}, nil
}

func TResolveAnotherKosmoStruct() (TKosmoStruct, error) {
	return TKosmoStruct{
		Name:   "Test",
		Number: 1,
	}, nil
}

func TestType(t *testing.T) {
	Convey("Given a structure", t, func() {
		schema := Type(TKosmoStruct{})
		Convey("It will reflect the type information of the given structure", func() {
			Convey("And will set the resolverType to the reflected graphql.Type", func() {
				infos := reflectTypeInformations(TKosmoStruct{})
				So(schema.resolverType, ShouldResemble, infos.typ)
			})
		})
	})
}

func TestQuery(t *testing.T) {
	Convey("Called on a GraphQLSchema, given a function", t, func() {
		schema := Type(TKosmoStruct{})
		schema.Query(TResolveKosmoStruct)
		Convey("The root 'Query' should be added", func() {
			Convey("And the given function should be used to resolve the previously given struct", func() {
				field := schema.query.Fields.(graphql.Fields)["TResolveKosmoStruct"]
				infos := reflectFunctionInformations(TResolveKosmoStruct)
				So(field.Args, ShouldResemble, infos.args)
				So(field.Type, ShouldResemble, schema.resolverType)
				So(field.Description, ShouldEqual, infos.description)
			})
		})
	})
}

func TestMutations(t *testing.T) {
	Convey("Called on a GraphQLSchema, given a function", t, func() {
		schema := Type(TKosmoStruct{})
		schema.Query(TResolveKosmoStruct)
		schema.Mutations(TResolveKosmoStruct)
		Convey("The root 'Mutation' should be added", func() {
			infos := reflectFunctionInformations(TResolveKosmoStruct)
			field := schema.mutations.Fields.(graphql.Fields)["TResolveKosmoStruct"]
			So(field.Args, ShouldResemble, infos.args)
			So(field.Description, ShouldEqual, infos.description)
			So(field.Type, ShouldResemble, schema.resolverType)
		})
	})
}

func TestSchemas(t *testing.T) {
	svc := Service{
		GraphQLConfig: GraphQLConfig{
			RemoveResolverPrefixes: true,
			ResolverPrefixes:       []string{"T"},
		},
	}
	schema1 := Type(TKosmoStruct{}).Query(TResolveKosmoStruct).Mutations(TResolveAnotherKosmoStruct)
	schema2 := Type(TAnotherKosmoStruct{}).Query(TResolveAnotherKosmoStruct).Mutations(TResolveKosmoStruct)
	Convey("Called on a Service, given no GraphQLSchema", t, func() {
		svcCopy1 := svc.Schemas()
		So(svcCopy1.graphQL.Query, ShouldBeNil)
		So(svcCopy1.graphQL.Mutation, ShouldBeNil)
	})
	Convey("Called on a Service, given one or more GraphQLSchema", t, func() {
		svcCopy2 := svc.Schemas(schema1, schema2)
		So(svcCopy2.graphQL.Query, ShouldNotBeNil)
		So(svcCopy2.graphQL.Mutation, ShouldNotBeNil)
	})
}

func TestServer(t *testing.T) {
	svc := Service{}
	schema1 := Type(TKosmoStruct{}).Query(TResolveKosmoStruct).Mutations(TResolveAnotherKosmoStruct)
	schema2 := Type(TAnotherKosmoStruct{}).Query(TResolveAnotherKosmoStruct).Mutations(TResolveKosmoStruct)
	Convey("Called on a Service", t, func() {
		Convey("In case the schema is not error free", func() {
			Convey("It should panic with the returned error message", func() {
				So(func() {
					svc.Schemas().Server()
				}, ShouldPanic)
			})
		})
		Convey("In case the schema is error free", func() {
			Convey("It should return a http.Server", func() {
				So(svc.Schemas(schema1, schema2).Server(), ShouldNotBeEmpty)
			})
		})
	})
}
