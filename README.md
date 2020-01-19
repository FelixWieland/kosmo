# kosmo

<img src="https://i.ibb.co/MspV6Mh/logo.png" align="right"
     title="Kosmo logo" width="120">

Kosmo is a microservice library whose goal is to offer GraphQL interfaces through native go types. As a GraphQL implementation it uses the "github.com/graphql-go/graphql" library. Because the GraphQL schema is created by native types and the resolver methods allow typed results to be returned, the service created can not only be used as a web interface, but can also be easily imported and used by other Go programs.

# Installation
```sh
go get github.com/FelixWieland/kosmo
```

# Usage

Here is a quick example to get you started:

```go
package main

import (
	"errors"

	"github.com/FelixWieland/kosmo"
)

//Passenger holds the data of a passenger
type Passenger struct {
	ID   int
	Name string `description:"Forename and Surname"`
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

//GetPassenger returns a Passenger
func GetPassenger(args ResolvePassengerArgs) (Passenger, error) {
	if args.ID == 0 {
		return Passenger{}, errors.New("Es ist ein fehler aufgetreten")
	}
	return Passenger{
		ID:   args.ID,
		Name: "Max",
		Seat: 1,
	}, nil
}

//GetPassengers returns multiple Passengers
func GetPassengers(args ResolvePassengersArgs) (Passengers, error) {
	return Passengers{
		Passenger{
			Name: "Max",
			Seat: 1,
		},
	}, nil
}

func main() {
	service := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Port: ":8080",
		},
		GraphQLConfig: kosmo.GraphQLConfig{
			ReplaceResolverPrefixes: true,
			ResolverPrefixes:        []string{"Get"},
		},
	}
	passenger := kosmo.Type(Passenger{}).Query(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Query(GetPassengers)

	service.Schemas(passenger, passengers).Server().ListenAndServe()
}

```
