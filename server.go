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

	rawHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: config.Playground,
	})

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
