# kosmo
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/FelixWieland/kosmo)
[![Actions Status](https://github.com/FelixWieland/kosmo/workflows/Test/badge.svg)](https://github.com/FelixWieland/kosmo/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/FelixWieland/kosmo)](https://goreportcard.com/report/github.com/FelixWieland/kosmo)
[![Coverage Status](https://coveralls.io/repos/github/FelixWieland/kosmo/badge.svg?branch=master&service=github)](https://coveralls.io/github/FelixWieland/kosmo?branch=master)

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
	passenger := kosmo.Type(Passenger{}).Queries(GetPassenger)
	passengers := kosmo.Type(Passengers{}).Queries(GetPassengers)

	service.Schemas(passenger, passengers).Server().ListenAndServe()
}
```

# Requests

Request your service by visiting "http://localhost:8080/" in your Browser and query:
```graphql
query {
	GetPassengers {
		Name
	}
}
```

With cURL:
```bash
curl --location --request POST 'http://localhost:8080/' \
	--header 'Content-Type: application/json' \
	--data-raw '{"query":"\nquery{\GetPassengers{\nName\n}\n}","variables":{}}'
```

With JS-Fetch:
```js
var myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");

var graphql = JSON.stringify({
  query: "\nquery{\GetPassengers {\nName\n}\n}",
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
			ResolverPrefixes: []string{},			// Prefixes that should be removed
		},
	}
```
