package kosmo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseTagConfig(t *testing.T) {
	Convey("Given a string", t, func() {
		config1 := parseTagConfig("require,otherStuff,   stuff, ignore")
		Convey("It should parse the string to a TagConfig an ignore all values that are not supported", func() {
			So(config1.Ignore, ShouldBeTrue)
			So(config1.Require, ShouldBeTrue)
		})
		config2 := parseTagConfig("")
		Convey("It should also work if the string is empty", func() {
			So(config2.Ignore, ShouldBeFalse)
			So(config2.Require, ShouldBeFalse)
		})
		config3 := parseTagConfig("require")
		Convey("And if only one viable flag is included", func() {
			So(config3.Ignore, ShouldBeFalse)
			So(config3.Require, ShouldBeTrue)
		})
		config4 := parseTagConfig("required")
		Convey("Should also work if require is provided as simple past", func() {
			So(config4.Ignore, ShouldBeFalse)
			So(config4.Require, ShouldBeTrue)
		})
		config5 := parseTagConfig("ignored")
		Convey("Should also work if ignore is provided as simple past", func() {
			So(config5.Ignore, ShouldBeTrue)
		})
	})
}
