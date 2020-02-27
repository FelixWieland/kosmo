package main

import (
	"github.com/FelixWieland/kosmo"
)

//ToIgnore should be ignored
type ToIgnore struct {
	field uint64 //this should fail
}

//Root -
type Root struct {
	ToIgnore `kosmo:"ignore"`
	ID       int
	Name     string
	Child    Child
}

//Child -
type Child struct {
	ID   int
	Name string
}

//GetRoot -
func GetRoot() (Root, error) {
	return Root{
		ID:   1,
		Name: "Root",
	}, nil
}

//GetAnotherRoot -
func GetAnotherRoot() (Root, error) {
	return Root{
		ID:   1,
		Name: "Another",
	}, nil
}

//GetChild -
func GetChild() (Child, error) {
	return Child{
		ID:   1,
		Name: "Child",
	}, nil
}

//MutateSomething tests a mutation
func MutateSomething() (Child, error) {
	return Child{
		ID:   1,
		Name: "Mutated",
	}, nil
}

//MutateAnother tests a mutation
func MutateAnother() (Root, error) {
	return Root{
		ID:   1,
		Name: "Mutated",
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
	root := kosmo.Type(Root{}).Queries(GetRoot, GetAnotherRoot).Mutations(MutateAnother)
	child := kosmo.Type(Child{}).Queries(GetChild)

	service.Schemas(root, child).Server().ListenAndServe()
}
