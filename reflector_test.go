package kosmo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"

	"github.com/kr/pretty"
)

type ActivityItem struct {
	ActivityItemID uint `gorm:"primary_key"`
	Variant        string
	Label          string
}

type ResolveActivityItemArguments struct {
	Name  string
	Index int
}

type ActivityItems []ActivityItem

func GetActivityItem(args ResolveActivityItemArguments) (ActivityItem, error) {
	return ActivityItem{
		Label: args.Name,
	}, nil
}

func prprint(val interface{}) {
	fmt.Printf("%# v", pretty.Formatter(val))
}

func TestReflectGormStruct(t *testing.T) {
	graph := structToGraph(ActivityItem{})
	if graph == nil {
		t.Fail()
	}
}

func TestReflectArgsFromResolver(t *testing.T) {
	args := functionToConfigArguments(reflect.ValueOf(GetActivityItem))
	if args == nil {
		t.Fail()
	}
	// prprint(args)
}

func TestResolverFactory(t *testing.T) {
	resolver := functionToResolver(reflect.ValueOf(GetActivityItem))
	args := make(map[string]interface{})
	args["Name"] = "test1"
	args["Index"] = 1
	value, err := resolver(graphql.ResolveParams{
		Args: args,
	})

	if err != nil {
		t.Fail()
	}

	structured := value.(ActivityItem)

	if structured.Label != "test1" {
		t.Fail()
	}
}

func TestReflectArray(t *testing.T) {
	graph := sliceToGraph(ActivityItems{})
	if graph == nil {
		t.Fail()
	}
}

type Product struct {
	ID    int
	Name  string
	Info  string
	Price int
}

type ProductList []Product

func TestReflectNativeField(t *testing.T) {
	f := nativeFieldToGraphQL(reflect.TypeOf(Product{}).Field(0))
	fa := graphql.Field{
		Type: graphql.Int,
	}

	if !reflect.DeepEqual(f, fa) {
		t.Fail()
	}
}

func TestReflectNativeFields(t *testing.T) {
	oconf := graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Info": &graphql.Field{
				Type: graphql.String,
			},
			"Price": &graphql.Field{
				Type: graphql.Int,
			},
		},
	}
	roconf := structToGraphConfig(Product{})

	t1 := getType(oconf.Fields)
	t2 := getType(roconf.Fields)

	if t1 != t2 {
		t.Fail()
	}

	if !reflect.DeepEqual(oconf.Fields, roconf.Fields) {
		t.Fail()
	}
}

func TestGrapQlTypeConfigReflection(t *testing.T) {
	oconf := graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Info": &graphql.Field{
				Type: graphql.String,
			},
			"Price": &graphql.Field{
				Type: graphql.Int,
			},
		},
	}
	roconf := structToGraphConfig(Product{})

	if !reflect.DeepEqual(oconf, roconf) {
		t.Fail()
	}
}

func TestGraphQLTypeReflection(t *testing.T) {
	productType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Product",
			Fields: graphql.Fields{
				"ID": &graphql.Field{
					Type: graphql.Int,
				},
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Info": &graphql.Field{
					Type: graphql.String,
				},
				"Price": &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)

	reflectedProductType := structToGraph(Product{})

	if !reflect.DeepEqual(reflectedProductType, productType) {
		t.Fail()
	}

	productTypeList := graphql.NewList(productType)

	config := sliceToGraphConfig(ProductList{})
	obj := graphqlObjectCache.Read(config.Name, gqlObjFallbackFactory(config)).(*graphql.Object)
	reflectedProductTypeList := graphql.NewList(obj)

	if !reflect.DeepEqual(productTypeList, reflectedProductTypeList) {
		t.Fail()
	}

}

type ResolverArgs struct {
	ID int
}

func ResolverArgsTestFn(args ResolverArgs) (interface{}, error) {
	return nil, nil
}

func TestFieldConfigArgumentReflection(t *testing.T) {
	args := graphql.FieldConfigArgument{
		"ID": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	}
	reflectedArgs := functionToConfigArguments(reflect.ValueOf(ResolverArgsTestFn))

	if !reflect.DeepEqual(args, reflectedArgs) {
		t.Fail()
	}
}

// BENCHMARKS

func BenchmarkResolverFactory(b *testing.B) {
	resolver := functionToResolver(reflect.ValueOf(ActivityItem{}))
	args := make(map[string]interface{})
	args["Name"] = "test1"
	args["Index"] = 1
	for n := 0; n < b.N; n++ {
		resolver(graphql.ResolveParams{
			Args: args,
		})
	}
}

func BenchmarkNonReflectedResolver(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetActivityItem(ResolveActivityItemArguments{
			"test2",
			1,
		})
	}
}

type InnerStructure struct {
	Number int
}

type NestedStructure struct {
	Name   string
	Nested InnerStructure
}

func TestNestedStructs(t *testing.T) {
	tInfos := reflectTypeInformations(NestedStructure{})
	_ = tInfos
}
