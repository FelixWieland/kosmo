package kosmo

import (
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

func getType(genVar interface{}) string {
	return reflect.ValueOf(genVar).Type().Name()
}

func rewriteFieldNamesToType(fields graphql.Fields) graphql.Fields {
	newFields := graphql.Fields{}
	for _, value := range fields {

		name := value.Type.Name()
		switch reflect.TypeOf(value.Type).Elem().Name() {
		case "List":
			newFields[name+"s"] = value //pluralize if list
		default:
			newFields[name] = value
		}
	}
	return newFields
}

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
