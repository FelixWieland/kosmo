package kosmo

import (
	"reflect"
	"runtime"
	"strings"
)

func runtimeFunctionName(function reflect.Value) string {
	nameWithPackage := runtime.FuncForPC(function.Pointer()).Name()
	parts := strings.SplitAfter(nameWithPackage, ".")
	return parts[len(parts)-1]
}

func reflectStructInformations(structType reflect.Type) (string, []reflect.StructField) {
	fields := []reflect.StructField{}
	for i := 0; i < structType.NumField(); i++ {
		config := parseTagConfig(structType.Field(i).Tag.Get("kosmo"))
		if config.Ignore {
			continue
		}
		fields = append(fields, structType.Field(i))
	}
	return structType.Name(), fields
}

func reflectArgumentFromResolverFunction(fn reflect.Value) (reflect.Type, bool) {
	if fn.Type().NumIn() == 0 {
		return nil, false
	}
	return fn.Type().In(0), true
}

func reflectTypeKind(value interface{}) string {
	return reflect.TypeOf(value).Kind().String()
}

func createFunctionStructArgumentFromMap(argumentType reflect.Type, argumentMap map[string]interface{}) reflect.Value {
	raw := reflect.New(argumentType).Elem()
	for key, field := range argumentMap {
		raw.FieldByName(key).Set(reflect.ValueOf(field))
	}
	return raw
}

func getType(genVar interface{}) string {
	return reflect.ValueOf(genVar).Type().Name()
}
