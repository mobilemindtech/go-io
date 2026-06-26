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
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
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

func NewOf[T any]() T {
	typOf := reflect.TypeFor[T]()

	// 1. Se for um Ponteiro (ex: *MyStruct ou *[]MyStruct)
	if typOf.Kind() == reflect.Pointer {
		elemType := typOf.Elem()

		// Se o ponteiro for para um Slice (ex: *[]MyStruct)
		if elemType.Kind() == reflect.Slice {
			emptySlice := reflect.MakeSlice(elemType, 0, 0)
			ptr := reflect.New(elemType)
			ptr.Elem().Set(emptySlice)
			return ptr.Interface().(T)
		}

		// Ponteiro para Struct normal (ex: *MyStruct)
		val := reflect.New(elemType).Interface()
		return val.(T)
	}

	// 2. Se for um Slice direto (ex: []MyStruct)
	if typOf.Kind() == reflect.Slice {
		val := reflect.MakeSlice(typOf, 0, 0).Interface()
		return val.(T)
	}

	// 3. Se for um Struct normal (ex: MyStruct)
	val := reflect.New(typOf).Elem().Interface().(T)
	return val
}
