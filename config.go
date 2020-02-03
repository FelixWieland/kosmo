package kosmo

import (
	"strings"
)

// TagConfig holds all flag configurations that can be made on a struct field
type TagConfig struct {
	Require bool
	Ignore  bool
}

func parseTagConfig(tagConfig string) TagConfig {
	flags := strings.Split(strings.ToLower(strings.ReplaceAll(tagConfig, " ", "")), ",")
	config := TagConfig{}
	for _, flag := range flags {
		switch flag {
		case "require":
			config.Require = true
		case "required":
			config.Require = true
		case "ignore":
			config.Ignore = true
		case "ignored":
			config.Ignore = true
		default:
		}
	}
	return config
}
