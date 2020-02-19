package kosmo

import (
	"net/http"

	"github.com/NYTimes/gziphandler"

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

	if config.Gzip {
		gzippedHandler := gziphandler.GzipHandler(rawHandler)
		mux.HandleFunc(config.APIBase, gzippedHandler.ServeHTTP)
	} else {
		mux.HandleFunc(config.APIBase, rawHandler.ServeHTTP)
	}

	return &server
}

func rawHandler(playground bool, schema graphql.Schema) *handler.Handler {
	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: playground,
	})
}
