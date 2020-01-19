package kosmo

import (
	"strings"

	"github.com/graphql-go/graphql"
)

func replaceResolverPrefixes(prefixes []string, fields graphql.Fields) graphql.Fields {
	newFields := graphql.Fields{}
	for key, value := range fields {
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				newFields[strings.Replace(key, prefix, "", 1)] = value
			}
		}

	}
	return newFields
}

func validateResolverArgument(resolverArg interface{}) (interface{}, string) {
	switch resolverArg.(type) {
	case Describer:
		return resolverArg.(Describer).Value, resolverArg.(Describer).Description
	default:
		return resolverArg, ""
	}
}

func gqlObjFallbackFactory(conf graphql.ObjectConfig) func(SetCache) {
	return func(set SetCache) {
		set(graphql.NewObject(conf))
	}
}
