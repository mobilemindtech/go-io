package types

import "github.com/mobilemindtec/go-io/result"

type IResourceIO interface {
	GetVarName() string
	Open() *result.Result[any]
	Close() *result.Result[any]
}

type ResourceIO[T any] struct {
	OpenFn  func() *result.Result[T]
	CloseFn func() *result.Result[T]
	VarName string
}

func (this *ResourceIO[T]) GetVarName() string {
	return this.VarName
}
func (this *ResourceIO[T]) Open() *result.Result[any] {
	return this.OpenFn().ToResultOfAny()
}
func (this *ResourceIO[T]) Close() *result.Result[any] {
	return this.CloseFn().ToResultOfAny()
}
