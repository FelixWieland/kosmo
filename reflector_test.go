package kosmo

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

type ActivityItem struct {
	ActivityItemID uint `gorm:"primary_key"`
	Variant        string
	Label          string
}

func prprint(val interface{}) {
	fmt.Printf("%# v", pretty.Formatter(val))
}

func TestReflectGormStruct(t *testing.T) {
	reflected := reflectGormType(ActivityItem{})
	graph := structToGraph(ActivityItem{})
	prprint(reflected)
	prprint(graph)
}
