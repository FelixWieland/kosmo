package kosmo

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TReflectStructInformationsStruct struct {
	Field1 string
	Field2 int
	Field3 bool
}

func TReflectArgumentFromResolverFunctionWithArgs(args TReflectStructInformationsStruct) {
}

func TestRuntimeFunctionName(t *testing.T) {
	Convey("Given a function", t, func() {
		Convey("The return value of #runtimeFunctionName should be its defined name", func() {
			functionName := runtimeFunctionName(reflect.ValueOf(TestRuntimeFunctionName))
			So(functionName, ShouldEqual, "TestRuntimeFunctionName")
		})
	})
}

func TestReflectStructInformations(t *testing.T) {
	Convey("Given a structure type", t, func() {
		name, fields := reflectStructInformations(reflect.TypeOf(TReflectStructInformationsStruct{}))

		Convey(`The first return value (name) should be the name of the type of the passed in type`, func() {
			So(name, ShouldEqual, "TReflectStructInformationsStruct")
		})
		Convey(`The second return value (fields) should have the lenght of the number of fields in that struct`, func() {
			So(len(fields), ShouldEqual, 3)
		})
		Convey(`Also should the second return value contains all Field names`, func() {
			So(fields[0].Name, ShouldEqual, "Field1")
			So(fields[1].Name, ShouldEqual, "Field2")
			So(fields[2].Name, ShouldEqual, "Field3")
		})
		Convey(`And the second return value contains all Field types`, func() {
			So(fields[0].Type.Name(), ShouldEqual, "string")
			So(fields[1].Type.Name(), ShouldEqual, "int")
			So(fields[2].Type.Name(), ShouldEqual, "bool")
		})
	})
}

func TestReflectArgumentFromResolverFunction(t *testing.T) {
	Convey("Given a function", t, func() {

		Convey("With no input parameters", func() {
			typ, hasArgs := reflectArgumentFromResolverFunction(reflect.ValueOf(func() {}))

			Convey("The first return value should be nil", func() {
				So(typ, ShouldEqual, nil)
			})
			Convey("The second return should be false (indicates no args)", func() {
				So(hasArgs, ShouldEqual, false)
			})
		})

		Convey("With some input parameters", func() {
			typ, hasArgs := reflectArgumentFromResolverFunction(reflect.ValueOf(TReflectArgumentFromResolverFunctionWithArgs))

			Convey("The first return value should be a valid type", func() {
				So(typ.Name(), ShouldEqual, "TReflectStructInformationsStruct")
			})
			Convey("The second return value should be true (indicates has args)", func() {
				So(hasArgs, ShouldEqual, true)
			})
		})
	})
}

func TestReflectTypeKind(t *testing.T) {
	Convey("Given a struct", t, func() {
		Convey("The return should be 'struct'", func() {
			So(reflectTypeKind(struct{}{}), ShouldEqual, "struct")
		})
	})
	Convey("Given a slice", t, func() {
		Convey("The return should be 'slice'", func() {
			So(reflectTypeKind([]int{}), ShouldEqual, "slice")
		})
	})
}

func TestCreateFunctionStructArgumentFromMap(t *testing.T) {
	Convey("Given a map and a reflect.Type{struct}", t, func() {
		strInfMap := make(map[string]interface{})

		strInfMap["Field1"] = "1"
		strInfMap["Field2"] = 2
		strInfMap["Field3"] = true

		filledStructRaw := createFunctionStructArgumentFromMap(reflect.TypeOf(TReflectStructInformationsStruct{}), strInfMap)
		filledStruct := filledStructRaw.Interface().(TReflectStructInformationsStruct)

		Convey("The return should have filled the struct", func() {
			So(filledStruct.Field1, ShouldEqual, "1")
			So(filledStruct.Field2, ShouldEqual, 2)
			So(filledStruct.Field3, ShouldEqual, true)
		})
	})
}

func TestGetType(t *testing.T) {
	Convey("Given a value", t, func() {
		typ1 := getType(TReflectStructInformationsStruct{})
		typ2 := getType(1)
		typ3 := getType("...")

		Convey("The return should be the name of its type", func() {
			So(typ1, ShouldEqual, "TReflectStructInformationsStruct")
			So(typ2, ShouldEqual, "int")
			So(typ3, ShouldEqual, "string")
		})
	})
}
