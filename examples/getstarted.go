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
func GetPassengers() (Passengers, error) {
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
			Playground: true,
			Port:       ":8080",
		},
		GraphQLConfig: kosmo.GraphQLConfig{
			RemoveResolverPrefixes: true,
			ResolverPrefixes:       []string{"Get"},
		},
	}
	passenger := kosmo.Type(Passenger{}).Query(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Query(GetPassengers)

	service.Schemas(passenger, passengers).Server().ListenAndServe()
}
