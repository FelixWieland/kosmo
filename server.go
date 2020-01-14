package kosmo

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

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

func startServer(config HTTPConfig, schema graphql.Schema) error {

	if config.APIBase == "" {
		config.APIBase = "/api"
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle(config.APIBase, h)
	http.ListenAndServe(config.Port, nil)
	return nil
	// http.HandleFunc(config.APIBase, func(w http.ResponseWriter, r *http.Request) {
	// 	result := executeQuery(r.URL.Query().Get("query"), schema)
	// 	json.NewEncoder(w).Encode(result)
	// })

	// return http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
}
