package main

import (
	"github.com/FelixWieland/kosmo"
)

//Root -
type Root struct {
	ID    int
	Name  string
	Child Child
}

//Child -
type Child struct {
	ID   int
	Name string
}

func GetRoot() (Root, error) {
	return Root{
		ID:   1,
		Name: "Root",
	}, nil
}

func GetAnotherRoot() (Root, error) {
	return Root{
		ID:   1,
		Name: "Another",
	}, nil
}

func GetChild() (Child, error) {
	return Child{
		ID:   1,
		Name: "Child",
	}, nil
}

func main() {
	service := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Playground: true,
			Port:       ":8080",
		},
		GraphQLConfig: kosmo.GraphQLConfig{
			RemoveResolverPrefixes: true,
			ResolverPrefixes:       []string{"Get"},
		},
	}
	root := kosmo.Type(Root{}).Queries(GetRoot, GetAnotherRoot)
	child := kosmo.Type(Child{}).Queries(GetChild)

	service.Schemas(root, child).Server().ListenAndServe()
}
