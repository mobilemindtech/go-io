package option

import (
	"fmt"
	"github.com/mobilemindtec/go-io/util"
	"reflect"
)

type IOption interface {
	IsOption() bool
	GetValue() interface{}
	IsEmpty() bool
}

type Option[T any] struct {
	value T
}

func Of[T any](it T) *Option[T] {
	return &Option[T]{value: it}
}

func None[T any]() *Option[T] {
	return &Option[T]{}
}

func (this *Option[T]) Get() T {
	return this.OrNil()
}

func (this *Option[T]) OrNil() T {
	return this.value
}

func (this *Option[T]) GetOrElse(v T) T {
	if this.IsEmpty() {
		return v
	}
	return this.Get()
}

func (this *Option[T]) OrElse(f func() *Option[T]) *Option[T] {
	if this.IsEmpty() {
		return f()
	}
	return this
}

func (this *Option[T]) IfEmpty(f func()) *Option[T] {
	if this.Empty() {
		f()
	}
	return this
}

func (this *Option[T]) IfNonEmpty(f func(T)) *Option[T] {
	if this.NonEmpty() {
		f(this.value)
	}
	return this
}

func (this *Option[T]) Foreach(f func(T)) *Option[T] {
	if this.NonEmpty() {
		f(this.value)
	}
	return this
}

func (this *Option[T]) Debug() {
	typ := reflect.TypeOf(this)
	fmt.Println(fmt.Sprintf("<DEBUG>: %v[value=%v]", typ, this.value))
}

func (this *Option[T]) Empty() bool {
	return util.IsNil(this.value)
}

func (this *Option[T]) NonEmpty() bool {
	return !util.IsNil(this.value)
}

func (this *Option[T]) IsOption() bool {
	return true
}
func (this *Option[T]) GetValue() interface{} {
	return this.Get()
}
func (this *Option[T]) IsEmpty() bool {
	return this.Empty()
}

func Filter[T any](v *Option[T], f func(T) bool) *Option[T] {
	if v.NonEmpty() {
		if f(v.Get()) {
			return v
		}
	}
	return None[T]()
}

func (this *Option[T]) String() string {
	if this.NonEmpty() {
		return fmt.Sprintf("Some(%v)", this.Get())
	} else {
		return "None"
	}
}

func Map[T any, R any](v1 *Option[T], f func(T) R) *Option[R] {
	if v1.NonEmpty() {
		return Of[R](f(v1.Get()))
	}
	return None[R]()
}

func FlatMap[T any, R any](v1 *Option[T], f func(T) *Option[R]) *Option[R] {
	if v1.NonEmpty() {
		f(v1.Get())
	}
	return None[R]()
}
