package main

import "github.com/FelixWieland/kosmo"

//Test demo type
type Test struct {
	Name string
}

//GetTest resolves a Test
func GetTest(params struct{ Name string }) (Test, error) {
	return Test{
		Name: params.Name,
	}, nil
}

func main() {
	s := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Playground: true,
			Port:       ":8080",
		},
	}

	demoType := kosmo.Type(Test{}).Queries(GetTest)

	s.Schemas(demoType).Server().ListenAndServe()
}
