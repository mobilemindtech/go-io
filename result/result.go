package result

import (
	"errors"
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/util"
	"reflect"
)

type IResult interface {
	IsResult() bool
	GetValue() interface{}
	GetError() error
	HasError() bool
}

type Unit struct {
}

func NewUnit() *Unit {
	return &Unit{}
}

type Result[T any] struct {
	value        T
	err          error
	lazy         func() (T, error)
	evaluated    bool
	errorChannel interface{}
}

func Try[T any](f func() (T, error)) *Result[T] {
	v, e := f()
	return Make(v, e)
}

func Lazy[T any](f func() (T, error)) *Result[T] {
	return &Result[T]{lazy: f}
}

func Make[T any](val T, e error) *Result[T] {
	return &Result[T]{value: val, err: e, evaluated: true}
}

func OfValue[T any](val T) *Result[T] {
	return &Result[T]{value: val, evaluated: true}
}

func Cast[T any](val interface{}) *Result[T] {
	if v, ok := val.(T); ok {
		return OfValue(v)
	}
	var x T
	return OfError[T](errors.New(fmt.Sprintf("type cast error %v to %v",
		reflect.TypeOf(val), reflect.TypeOf(x))))
}

func OfNil[T any]() *Result[T] {
	return &Result[T]{evaluated: true}
}

func OfError[T any](err error) *Result[T] {
	return &Result[T]{err: err, evaluated: true}
}

func (this *Result[T]) Evaluate() *Result[T] {
	return Try[T](this.lazy)
}

func (this *Result[T]) checkEvaluated() {
	if !this.evaluated {
		panic("computation not evaluated")
	}
}

func (this *Result[T]) ToOption() *option.Option[T] {
	this.checkEvaluated()
	return option.Of[T](this.value)
}

func (this *Result[T]) ToResultOfAny() *Result[any] {
	this.checkEvaluated()
	return Make[any](this.OrNil(), this.Error())
}

func (this *Result[T]) OptionNonEmpty() bool {
	return this.ToOption().NonEmpty()
}

func (this *Result[T]) OptionEmpty() bool {
	return this.ToOption().Empty()
}

func (this *Result[T]) Error() error {
	this.checkEvaluated()
	return this.err
}

func (this *Result[T]) Get() T {
	return this.OrNil()
}

func (this *Result[T]) OrNil() T {
	this.checkEvaluated()
	return this.value
}

func (this *Result[T]) IfError(f func(error)) *Result[T] {
	this.checkEvaluated()
	if this.IsError() {
		f(this.err)
	}
	return this
}

func (this *Result[T]) IfOk(f func(T)) *Result[T] {
	this.checkEvaluated()
	if this.IsResult() {
		f(this.value)
	}
	return this
}

func (this *Result[T]) IsError() bool {
	this.checkEvaluated()
	return this.err != nil
}

func (this *Result[T]) IsOk() bool {
	this.checkEvaluated()
	return !this.IsError()
}

func (this *Result[T]) Debug() {
	this.checkEvaluated()
	typ := reflect.TypeOf(this)
	fmt.Println(fmt.Sprintf("<DEBUG>: %v[value=%v, error=%v]", typ, this.value, this.err))
}

func (this *Result[T]) IsResult() bool {
	return true
}

func (this *Result[T]) GetValue() interface{} {
	return this.Get()
}

func (this *Result[T]) GetError() error {
	return this.Error()
}

func (this *Result[T]) HasError() bool {
	return this.IsError()
}

func (this *Result[T]) String() string {
	if this.IsError() {
		return fmt.Sprintf("Failure(%v)", this.Error())
	} else {
		return fmt.Sprintf("Ok(%v)", this.Get())
	}
}

type ResultM[A any, S any] struct {
	result *Result[A]
}

func NewResultM[A any, S any](r *Result[A]) *ResultM[A, S] {
	return &ResultM[A, S]{result: r}
}

func (this *ResultM[A, S]) isSome() bool {
	return util.IsNotNil(this.result) && this.result.IsOk() && this.result.ToOption().NonEmpty()
}

func (this *ResultM[A, S]) Filter(f func(A) bool) *ResultM[A, S] {
	if this.isSome() {
		if f(this.result.Get()) {
			return this
		}
		return NewResultM[A, S](OfNil[A]())
	}
	return this
}

func (this *ResultM[A, S]) Map(f func(A) S) *Result[S] {
	if this.isSome() {
		return OfValue(f(this.result.Get()))
	}
	return OfError[S](this.result.GetError())
}

func (this *ResultM[A, S]) FlatMap(f func(A) *Result[S]) *Result[S] {
	if this.isSome() {
		return f(this.result.Get())
	}
	return OfError[S](this.result.GetError())
}

func Filter[T any](v *Result[T], f func(T) bool) *Result[T] {
	if v.ToOption().NonEmpty() {
		if f(v.Get()) {
			return v
		}
	}
	return v
}

func Map[T any, R any](v *Result[T], f func(T) R) *Result[R] {
	if v.ToOption().NonEmpty() {
		return OfValue(f(v.Get()))
	}
	return OfError[R](v.Error())
}

func FlatMap[T any, R any](v *Result[T], f func(T) *Result[R]) *Result[R] {
	if v.ToOption().NonEmpty() {
		return f(v.Get())
	}
	return OfError[R](v.Error())
}
