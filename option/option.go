package option

import (
	"fmt"
	"github.com/mobilemindtech/go-io/util"
	"reflect"
)

type IOption interface {
	IsOption() bool
	GetValue() interface{}
	IsEmpty() bool
}

type _Option[T any] interface {
	isSome() bool
	isNone() bool
	get() T
	String() string
}

type _Some[T any] struct {
	value T
}

func _newSome[T any](value T) *_Some[T] {
	return &_Some[T]{value: value}
}

func (this _Some[T]) isSome() bool {
	return true
}
func (this _Some[T]) isNone() bool {
	return false
}

func (this _Some[T]) get() T {
	return this.value
}

func (this _Some[T]) String() string {
	return fmt.Sprintf("Some(%v)", this.value)
}

type _None[T any] struct {
}

func _newNone[T any]() *_None[T] {
	return &_None[T]{}
}

func (this _None[T]) isSome() bool {
	return false
}
func (this _None[T]) isNone() bool {
	return true
}

func (this _None[T]) get() T {
	panic("invalid call Get of None")
}

func (this _None[T]) String() string {
	return "None"
}

type Option[T any] struct {
	value _Option[T]
}

func Of[T any](it T) *Option[T] {
	if util.IsNil(it) {
		return None[T]()
	}
	return Some(it)
}

func Some[T any](it T) *Option[T] {
	return &Option[T]{value: _newSome(it)}
}

func None[T any]() *Option[T] {
	return &Option[T]{value: _newNone[T]()}
}

func (this *Option[T]) Get() T {
	return this.value.get()
}

func (this *Option[T]) OrNil() T {
	if this.NonEmpty() {
		return this.Get()
	}
	var x T
	return x
}

func (this *Option[T]) Or(v T) T {
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

func (this *Option[T]) Filter(f func(T) bool) *Option[T] {
	if this.NonEmpty() {
		if f(this.Get()) {
			return this
		} else {
			None[T]()
		}
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
		f(this.Get())
	}
	return this
}

func (this *Option[T]) Resolve(fnone func(), fsome func(T)) *Option[T] {
	if this.IsNone() {
		fnone()
	} else {
		fsome(this.Get())
	}
	return this
}

func (this *Option[T]) Foreach(f func(T)) *Option[T] {
	if this.NonEmpty() {
		f(this.Get())
	}
	return this
}

func (this *Option[T]) Debug() {
	typ := reflect.TypeOf(this)
	fmt.Println(fmt.Sprintf("<DEBUG>: %v[value=%v]", typ, this.value.String()))
}

func (this *Option[T]) IsSome() bool {
	return this.value.isSome()
}

func (this *Option[T]) IsNone() bool {
	return this.value.isNone()
}

func (this *Option[T]) Empty() bool {
	return this.value.isNone()
}

func (this *Option[T]) NonEmpty() bool {
	return this.value.isSome()
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

func (this *Option[T]) String() string {
	return this.value.String()
}

func (this *Option[T]) OrPanic(msg string) T {
	if this.IsSome() {
		return this.Get()
	}
	panic(msg)
}

func Filter[T any](v *Option[T], f func(T) bool) *Option[T] {
	if v.NonEmpty() {
		if f(v.Get()) {
			return v
		}
	}
	return None[T]()
}

func Map[T any, R any](v1 *Option[T], f func(T) R) *Option[R] {
	if v1.NonEmpty() {
		return Of[R](f(v1.Get()))
	}
	return None[R]()
}

func OrValue[T any](v1 *Option[T], value T) T {
	if v1.NonEmpty() {
		return v1.Get()
	}
	return value
}
func Or[T any](v1 *Option[T], f func() T) T {
	if v1.NonEmpty() {
		return v1.Get()
	}
	return f()
}

func OrElse[T any](v1 *Option[T], f func() *Option[T]) *Option[T] {
	if v1.NonEmpty() {
		return v1
	}
	return f()
}

func MapMaybe[T any, R any](v1 T, f func(T) R) *Option[R] {
	v := Of(v1)
	if v.NonEmpty() {
		return Of[R](f(v.Get()))
	}
	return None[R]()
}

func FlatMap[T any, R any](v1 *Option[T], f func(T) *Option[R]) *Option[R] {
	if v1.NonEmpty() {
		f(v1.Get())
	}
	return None[R]()
}
func Unwrap[T any]() func(*Option[T]) T {
	return func(opt *Option[T]) T {
		return opt.Get()
	}
}

func IsSome[T any](opt *Option[T]) bool {
	return opt.NonEmpty()
}
