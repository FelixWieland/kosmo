package kosmo

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

type Item struct {
	Name        string
	Description string
}

type Items []Item

type ResolveItemArguments struct{ name string }

func GetItem(args ResolveItemArguments) (Item, error) {
	prprint("Test")
	return Item{}, nil
}

func GetItems(args ResolveItemArguments) (Items, error) {
	return Items{}, nil
}

func TestKosmo(t *testing.T) {
	service := Service{
		HTTPConfig: HTTPConfig{
			Port: ":8080",
		},
		GraphQLConfig: GraphQLConfig{},
	}
	item := Type(Item{}).Query(GetItem)
	items := Type(Items{}).Query(GetItems)

	s := service.Schemas(item, items)

	s.Server().Close()
}

type Inner struct {
	Field2 string
}

type InnerSlice []Inner

type Test struct {
	Feld       string `description:"TestField"`
	Inner      Inner  `description:"Test-Nested"`
	InnerSlice InnerSlice
}

type ResolveTestArgs struct {
	Name string
}

func prprint(val interface{}) {
	fmt.Printf("%# v", pretty.Formatter(val))
}

func GetTest(args ResolveTestArgs) (Test, error) {
	return Test{
		Feld: args.Name,
		Inner: Inner{
			Field2: "InnerName",
		},
	}, nil
}

func TestMinimalExample(t *testing.T) {

	service := Service{
		HTTPConfig: HTTPConfig{
			Port:       ":8080",
			Playground: true,
		},
	}

	test := Type(Test{}).Query(GetTest)
	service.Schemas(test).Server()
}

func TestReplaceResolverPrefixExample(t *testing.T) {
	service := Service{
		HTTPConfig: HTTPConfig{
			Port: ":8080",
		},
		GraphQLConfig: GraphQLConfig{
			RemoveResolverPrefixes: true,
			ResolverPrefixes:       []string{"Get"},
		},
	}

	test := Type(Test{}).Query(Describer{Value: GetTest, Description: "Returns a Test"})
	service.Schemas(test).Server().Close()
}

func ResolverWithEmptyArgs() (Test, error) {
	return Test{}, nil
}

func TestEmptyResolver(t *testing.T) {
	testType := Type(Test{}).Query(ResolverWithEmptyArgs)
	_ = testType
}
