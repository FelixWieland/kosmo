package kosmo

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func muxServer(config HTTPConfig, schema graphql.Schema) *http.Server {
	if config.APIBase == "" {
		config.APIBase = "/"
	}

	rawHandler := rawHandler(config.Playground, schema)

	mux := http.NewServeMux()
	server := http.Server{Addr: config.Port, Handler: mux}

	mux.HandleFunc(config.APIBase, rawHandler.ServeHTTP)

	return &server
}

func rawHandler(playground bool, schema graphql.Schema) *handler.Handler {
	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: playground,
	})
}
