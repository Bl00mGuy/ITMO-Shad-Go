//go:build !solution

package illegal

import "unsafe"

func GetBytesPointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&b)
}

func ConvertPointerToString(ptr unsafe.Pointer) string {
	return *(*string)(ptr)
}

func StringFromBytes(b []byte) string {
	ptr := GetBytesPointer(b)
	return ConvertPointerToString(ptr)
}
