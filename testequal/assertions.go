//go:build !solution

package testequal

import (
	"bytes"
	"fmt"
	"reflect"
)

func Equal(expected, actual any) bool {
	if expected == nil || actual == nil {
		return checkNil(expected, actual)
	}

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}

	switch e := expected.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return checkPrimitiveTypes(expected, actual)
	case struct{}:
		return false
	case map[string]string:
		return checkMapStringString(e, actual)
	case []byte:
		return checkByteSlices(e, actual)
	default:
		return reflect.DeepEqual(expected, actual)
	}
}

func checkPrimitiveTypes(expected, actual any) bool {
	return expected == actual
}

func checkNil(expected, actual any) bool {
	return expected == nil && actual == nil
}

func checkMapStringString(expectedMap map[string]string, actual any) bool {
	mapA, ok := actual.(map[string]string)
	if !ok || len(expectedMap) != len(mapA) || len(mapA) == 0 {
		return false
	}

	for key, eVal := range expectedMap {
		aVal, ok := mapA[key]
		if !ok || eVal != aVal {
			return false
		}
	}
	return true
}

func checkByteSlices(expected, actual any) bool {
	bytesE, okE := expected.([]byte)
	bytesA, okA := actual.([]byte)
	if !okE || !okA {
		return false
	}

	if bytesE == nil || bytesA == nil {
		return bytesE == nil && bytesA == nil
	}

	return bytes.Equal(bytesE, bytesA)
}

func checkWithReflectDeepEqual(expected, actual any) bool {
	return reflect.DeepEqual(expected, actual)
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if Equal(expected, actual) {
		return true
	}

	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)

	return false
}

func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !Equal(expected, actual) {
		return true
	}

	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)

	return false
}

func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if Equal(expected, actual) {
		return
	}

	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)

	t.FailNow()
}

func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !Equal(expected, actual) {
		return
	}

	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)

	t.FailNow()
}
