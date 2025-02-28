//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func getFieldByName(obj interface{}, name string) reflect.Value {
	v := reflect.ValueOf(obj).Elem()
	return v.FieldByName(name)
}

func canSetField(field reflect.Value) bool {
	return field.IsValid() && field.CanSet()
}

func setFieldValueUnsafe(field reflect.Value, value interface{}) {
	ptr := unsafe.Pointer(field.UnsafeAddr())
	reflect.NewAt(field.Type(), ptr).Elem().Set(reflect.ValueOf(value))
}

func SetPrivateField(obj interface{}, name string, value interface{}) {
	field := getFieldByName(obj, name)
	if !canSetField(field) {
		setFieldValueUnsafe(field, value)
	}
}
