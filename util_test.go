package kosmo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReplaceResolverPrefixes(t *testing.T) {
	Convey("Given prefixes and graphql.Fields", t, func() {

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

	})
}

func TestMakeGraphQLObject(t *testing.T) {
	Convey("Given a graphql.ObjectConfig", t, func() {

	})
}
