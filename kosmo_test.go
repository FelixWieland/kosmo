package kosmo

import "testing"

type Item struct {
	Name        string
	Description string
}

type CreateItem Item

func (i Item) Resolve(args struct{ name string }) (interface{}, error) {
	prprint("Test")
	return Item{}, nil
}

func (i CreateItem) Resolve(args struct{ name string }) (interface{}, error) {
	prprint("Test2")
	return CreateItem{}, nil
}

func TestKosmo(t *testing.T) {
	item2 := CreateItem{
		Name: "test",
	}
	item2.Resolve(struct{ name string }{name: "test"})
}
