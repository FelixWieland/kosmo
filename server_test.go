package kosmo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TMuxServerSchemaStruct struct {
	Name string
}

func TResolveMuxServerSchema() (TMuxServerSchemaStruct, error) {
	return TMuxServerSchemaStruct{
		Name: "Test",
	}, nil
}

func TestMuxServer(t *testing.T) {
	svc := Service{
		HTTPConfig: HTTPConfig{
			Port: ":8080",
		},
	}
	tMux := Type(TMuxServerSchemaStruct{}).Query(TResolveMuxServerSchema)
	Convey("Given a HTTPConfig and a graphql.Schema", t, func() {
		server := svc.Schemas(tMux).Server()
		Convey("The returned server should have the port 8080", func() {
			So(server.Addr, ShouldEqual, ":8080")
		})
	})
}
