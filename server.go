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
