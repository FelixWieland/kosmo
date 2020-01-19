package kosmo

import "github.com/graphql-go/graphql"

func combineObjectConfig(configs ...graphql.ObjectConfig) graphql.ObjectConfig {
	combinedConfig := graphql.ObjectConfig{}
	fields := graphql.Fields{}
	for _, config := range configs {
		combinedConfig.Name = config.Name
		combinedConfig.Description = config.Description
		if config.Fields == nil {
			continue
		}
		for key, value := range config.Fields.(graphql.Fields) {
			if value == nil {
				continue
			}
			fields[key] = value
		}
	}
	combinedConfig.Fields = fields
	return combinedConfig
}

func makeGraphQLObject(objectConfig graphql.ObjectConfig) *graphql.Object {
	if objectConfig.Name == "" {
		return nil
	}
	obj := graphql.NewObject(objectConfig)
	return obj
}
