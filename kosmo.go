package kosmo

import "github.com/graphql-go/graphql"

// ResolveParams Params for Field.resolve()
type ResolveParams graphql.ResolveParams

// Field resolve the given type
type Field interface {
	resolve(struct{}) (interface{}, error)
}
