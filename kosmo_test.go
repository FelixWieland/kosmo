package kosmo

import "testing"

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
			Port: ":80",
		},
	}
	item := Type(Item{}).Query(GetItem)
	items := Type(Items{}).Query(GetItems)

	s := service.Schemas(item, items)

	prprint(s)

	s.Start()
}
