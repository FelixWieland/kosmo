package kosmo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	. "github.com/smartystreets/goconvey/convey"
)

type TForNestingStruct struct {
	Field1 string
}

type TForNestingSlice []TForNestingStruct

type TNativeFieldToGraphQLStruct struct {
	Field1 TForNestingStruct
	Field2 TForNestingSlice
	Field3 string `description:"Test"`
}

type TResolverArguments struct {
	Field1 string
	Field2 int
}

type TResolverArgumentsWithIngoredFields struct {
	Field1 string `kosmo:"ignore"`
	Field2 int
}

type TResolverArgumentsWithRequiredFields struct {
	Field1 string `kosmo:"require"`
	Field2 int
}

func TFunctionToResolverWithNoArgs() (TForNestingStruct, error) {
	return TForNestingStruct{
		Field1: "Test",
	}, nil
}

func TFunctionToResolverWithArgs(args TResolverArguments) (TForNestingStruct, error) {
	return TForNestingStruct{
		Field1: args.Field1,
	}, errors.New("Test")
}

func TFunctionToResolverWithArgsAndIgnoredFields(args TResolverArgumentsWithIngoredFields) (TForNestingStruct, error) {
	return TForNestingStruct{
		Field1: args.Field1,
	}, errors.New("Test")
}

func TFunctionToResolverWithArgsAndRequiredFields(args TResolverArgumentsWithRequiredFields) (TForNestingStruct, error) {
	return TForNestingStruct{
		Field1: args.Field1,
	}, errors.New("Test")
}

func TestReflectFunctionInformations(t *testing.T) {
	Convey("Given a function", t, func() {
		infos := reflectFunctionInformations(TFunctionToResolverWithArgs)
		Convey("It should assemble all informations needed for building a graphql resolver", func() {
			So(infos.metaInformations.name, ShouldEqual, "TFunctionToResolverWithArgs")
			So(infos.args, ShouldNotBeEmpty)
			So(infos.resolver, ShouldNotBeEmpty)
		})
	})
}

func TestReflectTypeInformations(t *testing.T) {
	Convey("Given a struct", t, func() {
		infos := reflectTypeInformations(TNativeFieldToGraphQLStruct{})
		Convey("It should assemble all informations needed for building a graphql schema", func() {
			So(infos.metaInformations.description, ShouldEqual, "")
			So(infos.typ, ShouldNotBeEmpty)
		})
	})
	Convey("Given a slice", t, func() {
		infos := reflectTypeInformations(TForNestingSlice{})
		Convey("It should assemble all informations needed for building a graphql schema", func() {
			So(infos.metaInformations.description, ShouldEqual, "")
			So(infos.typ, ShouldNotBeEmpty)
		})
	})
}

func TestFunctionToConfigArguments(t *testing.T) {
	Convey("Given a reflect.Value of a function", t, func() {
		Convey("In case the function has no args", func() {
			fieldConfigArgument := functionToConfigArguments(reflect.ValueOf(TFunctionToResolverWithNoArgs))
			Convey("It should return a empty graphql.FieldConfigArgument", func() {
				So(fieldConfigArgument, ShouldBeEmpty)
			})
		})
		Convey("In case the function has args", func() {
			fieldConfigArgument := functionToConfigArguments(reflect.ValueOf(TFunctionToResolverWithArgs))
			Convey("It should return a graphql.FieldConfigArgument with all the fields from its arguments type", func() {
				So(fieldConfigArgument, ShouldNotBeEmpty)
				So(fieldConfigArgument["Field1"], ShouldNotBeEmpty)
				So(fieldConfigArgument["Field2"], ShouldNotBeEmpty)
			})
		})
		Convey("In case some args in the function are marked as ignored", func() {
			fieldConfigArgument := functionToConfigArguments(reflect.ValueOf(TFunctionToResolverWithArgsAndIgnoredFields))
			Convey("The ignored fields should not be returned arguments", func() {
				So(fieldConfigArgument, ShouldNotBeEmpty)
				So(fieldConfigArgument["Field1"], ShouldBeNil)
				So(fieldConfigArgument["Field2"], ShouldNotBeEmpty)
			})
		})
		Convey("In case some args in the function are marked as required", func() {
			fieldConfigArgument := functionToConfigArguments(reflect.ValueOf(TFunctionToResolverWithArgsAndRequiredFields))
			Convey("The required function should be a non null type", func() {
				So(fieldConfigArgument, ShouldNotBeEmpty)
				So(fieldConfigArgument["Field1"], ShouldNotBeEmpty)
				So(fieldConfigArgument["Field2"], ShouldNotBeEmpty)
			})
		})
	})
}

func TestFunctionToResolver(t *testing.T) {
	Convey("Given a reflect.Value of a function", t, func() {
		Convey("In case the function has no args", func() {
			fn := functionToResolver(reflect.ValueOf(TFunctionToResolverWithNoArgs))
			Convey("It should call the passed in function and return it's results", func() {
				val, err := fn(graphql.ResolveParams{})
				So(val, ShouldResemble, TForNestingStruct{
					Field1: "Test",
				})
				So(err, ShouldBeNil)
			})
		})
		Convey("In case the function has args", func() {
			Convey("The returned function should map the graphql.ResolveParams to the passed in function and return its results", func() {
				args := make(map[string]interface{})
				args["Field1"] = "Test1"
				args["Field2"] = 2
				fn := functionToResolver(reflect.ValueOf(TFunctionToResolverWithArgs))
				val, err := fn(graphql.ResolveParams{
					Args: args,
				})
				So(val, ShouldResemble, TForNestingStruct{
					Field1: "Test1",
				})
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestStructToGraph(t *testing.T) {
	Convey("Given a struct", t, func() {
		obj := structToGraph(TForNestingStruct{})
		Convey("It should produce a object of its type", func() {
			So(obj.Name(), ShouldEqual, "TForNestingStruct")
		})
	})
}

func TestSliceToGraph(t *testing.T) {
	Convey("Given a slice", t, func() {
		list := sliceToGraph(TForNestingSlice{})
		Convey("It should produce a list of its underling type", func() {
			So(list.OfType.Name(), ShouldEqual, "TForNestingStruct")
		})
	})
}

func TestStructToGraphConfig(t *testing.T) {
	Convey("Given a struct", t, func() {
		objConf := structToGraphConfig(TForNestingStruct{})
		Convey("The name of the returned config should be the name of its type", func() {
			So(objConf.Name, ShouldEqual, "TForNestingStruct")
		})
	})
}

func TestSliceToGraphConfig(t *testing.T) {
	Convey("Given a slice", t, func() {
		objConf := sliceToGraphConfig(TForNestingSlice{})
		Convey("The name of the returned config should be the name from the underling type", func() {
			So(objConf.Name, ShouldEqual, "TForNestingStruct")
		})
	})
}

func TestBuildObjectConfigFromType(t *testing.T) {
	Convey("Given a reflect.Type of a struct", t, func() {
		objConf := buildObjectConfigFromType(reflect.TypeOf(TNativeFieldToGraphQLStruct{}))
		Convey("The name of the returned config should be the name from the type of the passed in struct", func() {
			So(objConf.Name, ShouldEqual, "TNativeFieldToGraphQLStruct")
		})
		Convey("The length of the returned fields should be the number of fields in given struct", func() {
			So(len(objConf.Fields.(graphql.Fields)), ShouldEqual, 3)
		})
	})
}

func TestNativeFieldsToGraphQLFields(t *testing.T) {
	Convey("Given reflected fields of a struct", t, func() {
		_, fields := reflectStructInformations(reflect.TypeOf(TNativeFieldToGraphQLStruct{}))
		gqlFields := nativeFieldsToGraphQLFields(fields)
		Convey("All reflected fields should be mapped to graphql.Fields", func() {
			Convey("The length of the returned fields should equal", func() {
				So(len(gqlFields), ShouldEqual, len(fields))
			})
			Convey("Names of the returned fields should be the name of the underling type", func() {
				So(gqlFields["Field1"].Type.Name(), ShouldEqual, "TForNestingStruct")
				So(gqlFields["Field2"].Type.Name(), ShouldEqual, "TForNestingStruct")
				So(gqlFields["Field3"].Type.Name(), ShouldEqual, "String")
			})
			Convey("In case a field has a 'description' tag, the description of the field should be the value of its tag", func() {
				So(gqlFields["Field3"].Description, ShouldEqual, "Test")
			})
			Convey("In case a field has no 'description' tag, the description of the field should be a empty string", func() {
				So(gqlFields["Field1"].Description, ShouldEqual, "")
			})
		})
	})
}

func TestNativeFieldToGraphQL(t *testing.T) {
	Convey("Given a reflected field of a struct", t, func() {
		_, fields := reflectStructInformations(reflect.TypeOf(TNativeFieldToGraphQLStruct{}))

		Convey("In case the input is a struct", func() {
			gqlField := nativeFieldToGraphQL(fields[0])
			Convey("The field type name should be the type of the struct", func() {
				So(gqlField.Type.Name(), ShouldEqual, "TForNestingStruct")
			})
		})
		Convey("In case the input is a slice", func() {
			gqlField := nativeFieldToGraphQL(fields[1])
			Convey("The field type name should be the name of the underling slice type", func() {
				So(gqlField.Type.Name(), ShouldEqual, "TForNestingStruct")
			})
		})
		Convey("In case the input is a primitive", func() {
			gqlField := nativeFieldToGraphQL(fields[2])
			Convey("The field type name should be the name of the mapped graphql.Type", func() {
				So(gqlField.Type.Name(), ShouldEqual, "String")
			})
		})
	})
}

func TestNativeTypeToGraphQL(t *testing.T) {
	Convey("Given a name of a primitive type", t, func() {
		var typeUInt uint
		var typeInt int
		var typeString string
		var typeFloat32 float32
		var typeFloat64 float64

		Convey("The corresponding graphQL type should be returned", func() {
			So(nativeTypeToGraphQL(getType(typeInt)), ShouldEqual, graphql.Int)
			So(nativeTypeToGraphQL(getType(typeUInt)), ShouldEqual, graphql.Int)
			So(nativeTypeToGraphQL(getType(typeString)), ShouldEqual, graphql.String)
			So(nativeTypeToGraphQL(getType(typeFloat32)), ShouldEqual, graphql.Float)
			So(nativeTypeToGraphQL(getType(typeFloat64)), ShouldEqual, graphql.Float)
		})
		Convey("The function should panic if the given type is not available", func() {
			So(func() {
				nativeTypeToGraphQL("garbage")
			}, ShouldPanic)
		})
	})
}
