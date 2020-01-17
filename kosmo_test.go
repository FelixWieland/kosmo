package kosmo

import (
	"testing"
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
		GraphQLConfig: GraphQLConfig{
			UseTypeAsQueryName: true,
		},
	}
	item := Type(Item{}).Query(GetItem)
	items := Type(Items{}).Query(GetItems)

	s := service.Schemas(item, items)

	s.Server().Close()
}

type Test struct {
	Feld string `description:"TestField"`
}

type ResolveTestArgs struct {
	Name string
}

func GetTest(args ResolveTestArgs) (Test, error) {
	return Test{
		Feld: args.Name,
	}, nil
}

func TestMinimalExample(t *testing.T) {

	service := Service{
		HTTPConfig: HTTPConfig{
			Port: ":8080",
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
			ReplaceResolverPrefixes: true,
			ResolverPrefixes:        []string{"Get"},
		},
	}

	test := Type(Test{}).Query(Describer{Value: GetTest, Description: "Returns a Test"})
	service.Schemas(test).Server().Close()
}
