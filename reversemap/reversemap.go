//go:build !solution

package reversemap

import (
	"reflect"
)

func ReverseMap(forward interface{}) interface{} {
	forwardValue := reflect.ValueOf(forward)
	keyType := forwardValue.Type().Key()
	valueType := forwardValue.Type().Elem()
	reversed := reflect.MakeMap(reflect.MapOf(valueType, keyType))

	for _, key := range forwardValue.MapKeys() {
		value := forwardValue.MapIndex(key)
		reversed.SetMapIndex(value, key)
	}

	return reversed.Interface()
}
