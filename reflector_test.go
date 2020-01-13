package kosmo

import (
	"fmt"
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
	// reflected := reflectGormType(ActivityItem{})
	graph := structToGraph(ActivityItem{})
	if graph == nil {
		t.Fail()
	}
	// prprint(reflected)
	// prprint(graph)
}

func TestReflectArgsFromResolver(t *testing.T) {
	args := reflectArgsFromResolver(reflectResolverFunction(GetActivityItem))
	if args == nil {
		t.Fail()
	}
	// prprint(args)
}

func TestResolverFactory(t *testing.T) {
	resolver := resolverFactory(reflectResolverFunction(GetActivityItem))
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

func BenchmarkResolverFactory(b *testing.B) {
	resolver := resolverFactory(reflectResolverFunction(ActivityItem{}))
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

func TestReflectArray(t *testing.T) {
	graph := sliceToGraph(ActivityItems{})
	if graph == nil {
		t.Fail()
	}
}
