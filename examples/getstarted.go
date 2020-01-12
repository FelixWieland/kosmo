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
