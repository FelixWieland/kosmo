package kosmo

import (
	"encoding/json"
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

func (a ActivityItem) Resolve(args struct {
	Name  string
	Index int
}) (interface{}, error) {
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
	args := reflectArgsFromResolver(reflectResolverMethod(ActivityItem{}))
	if args == nil {
		t.Fail()
	}
	// prprint(args)
}

func TestResolverFactory(t *testing.T) {
	resolver := resolverFactory(reflectResolverMethod(ActivityItem{}))
	args := make(map[string]interface{})
	args["name"] = "test1"
	args["index"] = 1
	value, _ := resolver(graphql.ResolveParams{
		Args: args,
	})
	prprint(value)
}

func BenchmarkResolverFactory(b *testing.B) {
	resolver := resolverFactory(reflectResolverMethod(ActivityItem{}))
	args := make(map[string]interface{})
	args["name"] = "test1"
	args["index"] = 1
	for n := 0; n < b.N; n++ {
		resolver(graphql.ResolveParams{
			Args: args,
		})
	}
}

func BenchmarkNonReflectedResolver(b *testing.B) {
	item := ActivityItem{}
	for n := 0; n < b.N; n++ {
		item.Resolve(struct {
			Name  string
			Index int
		}{
			"test2",
			1,
		})
	}
}

func TestMapstruct(t *testing.T) {
	args := make(map[string]interface{})
	args["label"] = "test"
	args["index"] = 1

	store := ActivityItem{}

	jsonbody, _ := json.Marshal(args)
	if err := json.Unmarshal(jsonbody, &store); err != nil {
		// do error check
	}
	// prprint(store)
}
