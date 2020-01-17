package kosmo

var graphqlObjectCache Cache

func init() {
	graphqlObjectCache = NewCache()
}
