package kosmo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCache(t *testing.T) {
	Convey("Should create a new cache struct and return its address", t, func() {
		c1 := NewCache()
		Convey("With a empty string-interface map", func() {
			So(c1, ShouldResemble, &cache{
				store: make(map[string]interface{}),
			})
		})
	})
}

func TestRead(t *testing.T) {
	c1 := NewCache()
	Convey("Given a key and a fallback function", t, func() {
		key := "key1"
		valueToSet := "value1"
		shouldBeUsed := func(setter SetCache) {
			setter(valueToSet)
		}
		shouldNeverBeUsed := func(setter SetCache) {
			setter("this should never be set")
		}
		Convey("In case the cache has no item that matches the key", func() {
			Convey("The value in the fallback function should be set to the cache and returned", func() {
				So(c1.Read(key, shouldBeUsed).(string), ShouldEqual, valueToSet)
				So(c1.Read(key, shouldNeverBeUsed).(string), ShouldEqual, valueToSet)
			})
		})
		Convey("In case the cache has a item that matches the key", func() {
			Convey("The value returned from Read should be the value stored in the cache", func() {
				So(c1.Read(key, shouldNeverBeUsed).(string), ShouldEqual, valueToSet)
			})
		})
	})
}
