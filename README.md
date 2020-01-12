# kosmo

<img src="https://i.ibb.co/MspV6Mh/logo.png" align="right"
     title="Kosmo logo" width="120">

Kosmo is a microservice framework whose goal is to offer GraphQL interfaces through native go types. As a GraphQL implementation it uses the "github.com/graphql-go/graphql" library. Because the GraphQL schema is created by native types and the resolver methods allow typed results to be returned, the service created can not only be used as a web interface, but can also be easily imported and used by other Go programs.

# Installation
```sh
go get github.com/FelixWieland/kosmo
```

# Usage

Here is a quick example to get you started:

```go
package examples

import "github.com/FelixWieland/kosmo"

//Passenger holds the data of a passenger
type Passenger struct {
	ID   int
	Name string
	Seat int
}

//Passengers holds multiple Passenger
type Passengers []Passenger

//ResolvePassengerArgs used to resolve a passenger
type ResolvePassengerArgs struct {
	ID int
}

//ResolvePassengersArgs used to resolve multiple passengers
type ResolvePassengersArgs struct{}

//Resolve returns a Passenger
func (p Passenger) Resolve(args ResolvePassengerArgs) (Passenger, error) {
	return Passenger{
		ID:   args.ID,
		Name: "Max",
		Seat: 1,
	}, nil
}

//Resolve returns multiple Passengers
func (ps Passengers) Resolve(args ResolvePassengersArgs) (Passengers, error) {
	return Passengers{}, nil
}

func main() {
	service := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Port: 80,
		},
	}
	service.Queries(Passenger{}, Passengers{}).Start()
}
```
