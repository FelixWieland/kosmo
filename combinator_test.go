package kosmo

import (
	"testing"

	"github.com/graphql-go/graphql"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCombineObjectConfig(t *testing.T) {
	Convey("Given multiple graphql.ObjectConfig's", t, func() {
		config1 := graphql.ObjectConfig{
			Name: "Name1",
			Fields: graphql.Fields{
				"Test1": &graphql.Field{
					Name: "Test1",
				},
				"Test2": &graphql.Field{
					Name: "Test2",
				},
			},
		}
		config2 := graphql.ObjectConfig{
			Name: "Name1",
			Fields: graphql.Fields{
				"Test3": &graphql.Field{
					Name: "Test3",
				},
				"Test4": &graphql.Field{
					Name: "Test4",
				},
				"Test5": &graphql.Field{
					Name: "Test5",
				},
			},
		}
		config3 := graphql.ObjectConfig{
			Name:   "Name1",
			Fields: nil,
		}
		config4 := graphql.ObjectConfig{
			Name: "Name1",
			Fields: graphql.Fields{
				"Test6": &graphql.Field{
					Name: "Test3",
				},
				"Test7": nil,
			},
		}
		combinedConfig := combineObjectConfig(config1, config2, config3, config4)
		Convey("The graphql.ObjectConfig returned should contain all the fields of the passed in configs", func() {
			fields := combinedConfig.Fields.(graphql.Fields)
			_, firstIsOk := fields["Test1"]
			_, secondIsOk := fields["Test2"]
			_, thirdIsOk := fields["Test3"]
			_, fourthIsOk := fields["Test4"]
			_, fifthIsOk := fields["Test5"]

			So(firstIsOk, ShouldBeTrue)
			So(secondIsOk, ShouldBeTrue)
			So(thirdIsOk, ShouldBeTrue)
			So(fourthIsOk, ShouldBeTrue)
			So(fifthIsOk, ShouldBeTrue)
		})
	})
}
