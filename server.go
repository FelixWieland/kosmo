package kosmo

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func server(config HTTPConfig, schema graphql.Schema) *http.Server {
	if config.APIBase == "" {
		config.APIBase = "/"
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: config.Playground,
	})

	mux := http.NewServeMux()
	server := http.Server{Addr: config.Port, Handler: mux}
	mux.HandleFunc(config.APIBase, h.ServeHTTP)

	return &server
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
