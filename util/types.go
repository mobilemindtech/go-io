package util

import (
	"fmt"
	"reflect"
)

func IsNotNil(i interface{}) bool {
	return !IsNil(i)
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

func CanNil(kind reflect.Kind) bool {
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}

func PanicCastType(label string, typOfA reflect.Type, typOfB reflect.Type) {
	panic(fmt.Sprintf("can't cast %v to %v on %v", typOfA, typOfB, label))
}
