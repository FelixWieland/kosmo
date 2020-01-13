package main

import (
	"fmt"

	"github.com/FelixWieland/kosmo"
)

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

//GetPassenger returns a Passenger
func GetPassenger(args ResolvePassengerArgs) (Passenger, error) {
	return Passenger{
		ID:   args.ID,
		Name: "Max",
		Seat: 1,
	}, nil
}

//GetPassengers returns multiple Passengers
func GetPassengers(args ResolvePassengersArgs) (Passengers, error) {
	return Passengers{}, nil
}

func main() {
	service := kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			Port: ":80",
		},
	}
	passenger := kosmo.Type(Passenger{}).Query(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Query(GetPassengers)

	fmt.Printf("%#v", passengers)

	service.Schemas(passenger, passengers).Start()
}
