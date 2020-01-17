package kosmo

// Cache interface
type Cache interface {
	Read(string, func(SetCache)) interface{}
}

type cache struct {
	latestFallback string
	store          map[string]interface{}
}

// SetCache used to set a cache
type SetCache func(interface{})

//NewCache - Creates a new Cache
func NewCache() Cache {
	return &cache{
		store: make(map[string]interface{}),
	}
}

// Read - Reads from Cache
func (c *cache) Read(key string, fallback func(setter SetCache)) interface{} {
	if value, ok := c.store[key]; ok {
		return value
	}
	setter := func(value interface{}) {
		c.store[key] = value
	}
	fallback(setter)
	return c.store[key]
}
