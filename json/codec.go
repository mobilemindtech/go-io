package json

import (
	"encoding/json"
	"github.com/mobilemindtec/go-io/result"
	"reflect"
)

type JsonEncoder[T any] struct {
}

func NewJsonEncoder[T any]() *JsonEncoder[T] {
	return &JsonEncoder[T]{}
}

func (this *JsonEncoder[T]) Encode(data T) *result.Result[[]byte] {
	return result.Try(func() ([]byte, error) {
		return json.Marshal(data)
	})
}

type JsonDecoder[T any] struct {
}

func NewJsonDecoder[T any]() *JsonDecoder[T] {
	return &JsonDecoder[T]{}
}

func (this *JsonDecoder[T]) Decode(data []byte) *result.Result[T] {
	return result.Try(func() (T, error) {
		typOf := reflect.TypeFor[T]()
		if typOf.Kind() == reflect.Pointer {
			typOf = typOf.Elem()
			val := reflect.New(typOf).Interface()
			return val.(T), json.Unmarshal(data, val)
		} else {
			val := reflect.New(typOf).Elem().Interface().(T)
			return val, json.Unmarshal(data, &val)
		}
	})
}

func (this *JsonDecoder[T]) DecodeTo(data []byte, entity T) *result.Result[T] {
	return result.Try(func() (T, error) {
		return entity, json.Unmarshal(data, entity)
	})
}

func Decode[T any](data []byte) *result.Result[T] {
	return NewJsonDecoder[T]().Decode(data)
}

func DecodeTo[T any](data []byte, entity T) *result.Result[T] {
	return NewJsonDecoder[T]().DecodeTo(data, entity)
}

func Encode[T any](entity T) *result.Result[[]byte] {
	return NewJsonEncoder[T]().Encode(entity)
}
