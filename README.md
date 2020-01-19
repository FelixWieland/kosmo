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
			RemoveResolverPrefixes: true,
			ResolverPrefixes:        []string{"Get"},
		},
	}
	passenger := kosmo.Type(Passenger{}).Query(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Query(GetPassengers)

	service.Schemas(passenger, passengers).Server().ListenAndServe()
}

```

# Requests

Request your service by visiting "http://localhost:8080/" in your Browser and query:
```graphql
query {
	Passengers {
		Name
	}
}
```

With cURL:
```bash
curl --location --request POST 'http://localhost:8080/' \
	--header 'Content-Type: application/json' \
	--data-raw '{"query":"\nquery{\nPassengers{\nName\n}\n}","variables":{}}'
```

With JS-Fetch:
```js
var myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");

var graphql = JSON.stringify({
  query: "\nquery{\nPassengers {\nName\n}\n}",
  variables: {}
})
var requestOptions = {
  method: 'POST',
  headers: myHeaders,
  body: graphql,
  redirect: 'follow'
};

fetch("http://localhost:8080/", requestOptions)
  .then(response => response.text())
  .then(result => console.log(result))
  .catch(error => console.log('error', error));
```


# Configurations

A kosmo Service can take the following configurations (the following values are default):

```go
kosmo.Service{
		HTTPConfig: kosmo.HTTPConfig{
			APIBase: "/", 				// Root of the endpoint
			Port: ":8080",				// Port of the service
			Playground: false,			// GraphIQL Playground
		},
		GraphQLConfig: kosmo.GraphQLConfig{
			RemoveResolverPrefixes: false,			// Removes the given prefixes from the resolver names 
			ResolverPrefixes:        []string{},	// Prefixes that should be removed
		},
	}
```
