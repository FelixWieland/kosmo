package kosmo

import (
	"testing"

	"github.com/graphql-go/graphql"

	. "github.com/smartystreets/goconvey/convey"
)

type TStructForObjectConfigTest struct {
	Name string
}

func TestReplaceResolverPrefixes(t *testing.T) {
	Convey("Given prefixes and graphql.Fields", t, func() {
		fields := graphql.Fields{
			"Pre1Name1": &graphql.Field{},
			"Pre2Name2": &graphql.Field{},
		}
		prefixes := []string{"Pre1", "Pre2"}
		newFields := replaceResolverPrefixes(prefixes, fields)
		Convey("The matching prefixes in the given fields should have been replaced with an empty string", func() {
			_, ok1 := newFields["Name1"]
			_, ok2 := newFields["Name2"]

			So(ok1, ShouldBeTrue)
			So(ok2, ShouldBeTrue)
		})
	})
}

func TestValidateResolverArgument(t *testing.T) {
	Convey("Given resolver arguments", t, func() {
		fn := func() {}
		Convey("In a Describer", func() {
			rslvr, descr := validateResolverArgument(Describer{Value: fn, Description: "Test"})
			Convey("It should return the resolver and the provided description", func() {
				So(rslvr, ShouldEqual, fn)
				So(descr, ShouldEqual, "Test")
			})
		})
		Convey("Raw", func() {
			rslvr, descr := validateResolverArgument(fn)
			Convey("It should return the resolver and a empty description", func() {
				So(rslvr, ShouldEqual, fn)
				So(descr, ShouldEqual, "")
			})
		})
	})
}

func TestGqlObjFallbackFactory(t *testing.T) {
	Convey("Given a graphql.ObjectConfig", t, func() {
		conf := structToGraphConfig(TStructForObjectConfigTest{})
		fn := gqlObjFallbackFactory(conf)
		Convey("It should return a cache setter function", func() {
			Convey("And the setter should receive a new graphql.Object", func() {
				obj := graphql.NewObject(conf)
				setter := func(value interface{}) {
					So(value, ShouldHaveSameTypeAs, obj)
					So(value, ShouldNotBeNil)
				}
				fn(setter)
			})
		})
	})
}

func TestMakeGraphQLObject(t *testing.T) {
	Convey("Given a graphql.ObjectConfig", t, func() {
		conf := structToGraphConfig(TStructForObjectConfigTest{})
		Convey("In case the config name is empty", func() {
			confCopy := conf
			confCopy.Name = ""
			obj := makeGraphQLObject(confCopy)
			Convey("It should return nil", func() {
				So(obj, ShouldBeNil)
			})
		})
		Convey("In case the config name is set", func() {
			obj := makeGraphQLObject(conf)
			Convey("It should return a valid graphql.Object", func() {
				So(obj, ShouldNotBeNil)
			})
		})
	})
}
