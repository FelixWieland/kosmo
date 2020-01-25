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

var demoPassengers Passengers = Passengers{
	Passenger{
		ID:   0,
		Name: "Max",
		Seat: 1,
	},
	Passenger{
		ID:   1,
		Name: "Mia",
		Seat: 2,
	},
}

//GetPassenger returns a Passenger
func GetPassenger(args ResolvePassengerArgs) (Passenger, error) {
	for _, passenger := range demoPassengers {
		if passenger.ID == args.ID {
			return passenger, nil
		}
	}
	return Passenger{}, errors.New("Passenger not found")
}

//GetPassengers returns multiple Passengers
func GetPassengers() (Passengers, error) {
	return demoPassengers, nil
}

func main() {
	service := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Playground: true,
			Port:       ":8080",
		},
	}
	passenger := kosmo.Type(Passenger{}).Query(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Query(GetPassengers)

	service.Schemas(passenger, passengers).Server().ListenAndServe()
}
