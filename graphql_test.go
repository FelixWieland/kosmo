package kosmo

import (
	"testing"

	"github.com/graphql-go/graphql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReflectFunctionInformations(t *testing.T) {
	Convey("Given a function", t, func() {

	})
}

func TestReflectTypeInformations(t *testing.T) {
	Convey("Given a struct", t, func() {

	})
}

func TestFunctionToConfigArguments(t *testing.T) {
	Convey("Given a reflect.Value of a function", t, func() {

	})
}

func TestFunctionToResolver(t *testing.T) {
	Convey("Given reflect.Value of a function", t, func() {

	})
}

func TestStructToGraph(t *testing.T) {
	Convey("Given a struct", t, func() {

	})
}

func TestSliceToGraph(t *testing.T) {
	Convey("Given a slice", t, func() {

	})
}

func TestStructToGraphConfig(t *testing.T) {
	Convey("Given a struct", t, func() {

	})
}

func TestSliceToGraphConfig(t *testing.T) {
	Convey("Given a slice", t, func() {

	})
}

func TestBuildObjectConfigFromType(t *testing.T) {
	Convey("Given a reflect.Type of a struct or slice", t, func() {

	})
}

func TestNativeFieldsToGraphQLFields(t *testing.T) {
	Convey("Given reflected fields of a struct", t, func() {

	})
}

func TestNativeFieldToGraphQL(t *testing.T) {
	Convey("Given a reflected field of a struct", t, func() {

	})
}

func TestNativeTypeToGraphQL(t *testing.T) {
	Convey("Given a name of a primitive type", t, func() {

		var typeUInt uint
		var typeInt int
		var typeString string
		var typeFloat32 float32
		var typeFloat64 float32

		Convey("The corresponding graphQL type should be returned", func() {
			So(nativeTypeToGraphQL(getType(typeInt)), ShouldEqual, graphql.Int)
			So(nativeTypeToGraphQL(getType(typeUInt)), ShouldEqual, graphql.Int)
			So(nativeTypeToGraphQL(getType(typeString)), ShouldEqual, graphql.String)
			So(nativeTypeToGraphQL(getType(typeFloat32)), ShouldEqual, graphql.Float)
			So(nativeTypeToGraphQL(getType(typeFloat64)), ShouldEqual, graphql.Float)
		})

		Convey("The function should panic if the given type is not available", func() {
			defer func() {
				var paniced bool
				if r := recover(); r != nil {
					paniced = true
				}
				So(paniced, ShouldEqual, true)
			}()
			nativeTypeToGraphQL("garbage")
		})
	})
}
