package kosmo

import "testing"

type Item struct {
	Name        string
	Description string
}

type Items []Item

type CreateItem Item

type ResolveItemArguments struct{ name string }

func (i Item) Resolve(args ResolveItemArguments) (interface{}, error) {
	prprint("Test")
	return Item{}, nil
}

func (i CreateItem) Resolve(args ResolveItemArguments) (interface{}, error) {
	return CreateItem{}, nil
}

func (is Items) Resolve(args ResolveItemArguments) (interface{}, error) {
	return Items{}, nil
}

func TestKosmo(t *testing.T) {
	item2 := CreateItem{
		Name: "test",
	}

	ksm := Service{}
	ksm.Queries(Item{}, Items{}).Mutations(CreateItem{}).Start()

	item2.Resolve(ResolveItemArguments{name: "test"})
}
