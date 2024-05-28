package result

import (
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
	Error() string
	ToResultOf() *Result[any]
	ToResultOfOption() *Result[*option.Option[any]]
}

type Unit struct {
}

func NewUnit() *Unit {
	return &Unit{}
}

type _Result[T any] interface {
	IsOk() bool
	IsFailure() bool
	IsEmpty() bool
	Get() T
}

type _Ok[T any] struct {
	value T
}

func _newOk[T any](value T) *_Ok[T]  { return &_Ok[T]{value} }
func (this *_Ok[T]) IsOk() bool      { return true }
func (this *_Ok[T]) IsFailure() bool { return false }
func (this *_Ok[T]) Get() T          { return this.value }

type _Failure struct {
	err error
}

func _newFailure(err error) *_Failure  { return &_Failure{err} }
func (this *_Failure) IsOk() bool      { return false }
func (this *_Failure) IsFailure() bool { return true }
func (this *_Failure) Get() error      { return this.err }

type Result[T any] struct {
	ok           *_Ok[T]
	failure      *_Failure
	lazy         func() (T, error)
	evaluated    bool
	errorChannel interface{}
}

func Try[T any](f func() (T, error)) *Result[T] {
	v, e := f()
	return Make(v, e)
}

func TryOption[T any](f func() (T, error)) *Result[*option.Option[T]] {
	v, e := f()
	return MakeOption(v, e)
}

func Lazy[T any](f func() (T, error)) *Result[T] {
	return &Result[T]{lazy: f}
}

func Make[T any](val T, e error) *Result[T] {
	if util.IsNotNil(e) {
		return OfError[T](e)
	}
	return OfValue[T](val)
}

func MakeOption[T any](val T, e error) *Result[*option.Option[T]] {
	if util.IsNotNil(e) {
		return OfError[*option.Option[T]](e)
	}
	return OfValue[*option.Option[T]](option.Of(val))
}

func TryMake[T any](r *Result[any]) *Result[T] {

	if r.IsError() {
		return OfError[T](r.Failure())
	}

	if r.IsOk() {
		if v, ok := r.Get().(T); ok {
			return OfValue(v)
		}

		var t T
		panic(fmt.Sprintf("can't cast received value %v to %v",
			reflect.TypeOf(r.Get()), reflect.TypeOf(t)))
	}

	panic(fmt.Sprintf("Can't create empty result of %v", r))
}

func OfValue[T any](val T) *Result[T] {
	return &Result[T]{ok: _newOk(val), evaluated: true}
}

/*
func OfNil[T any]() *Result[T] {
	return &Result[T]{evaluated: true}
}*/

func Cast[T any](val interface{}) *Result[T] {
	if v, ok := val.(T); ok {
		return OfValue(v)
	}
	var x T
	panic(fmt.Sprintf("type cast error %v to %v",
		reflect.TypeOf(val), reflect.TypeOf(x)))
}

func OfError[T any](err error) *Result[T] {
	return &Result[T]{failure: _newFailure(err), evaluated: true}
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
	if this.ok != nil {
		return option.Of[T](this.ok.Get())
	}
	return option.None[T]()
}

func (this *Result[T]) ToResultOf() *Result[any] {
	this.checkEvaluated()

	if this.IsError() {
		return OfError[any](this.Failure())
	} else if this.IsOk() {
		return OfValue[any](this.Get())
	}

	panic("Invalid empty result")
}

func (this *Result[T]) ToResultOfOption() *Result[*option.Option[any]] {
	this.checkEvaluated()

	if this.IsError() {
		return OfError[*option.Option[any]](this.Failure())
	} else if this.IsOk() {

		if util.IsNotNil(this.GetValue()) {
			if opt, ok := this.GetValue().(option.IOption); ok {
				if !opt.IsEmpty() {
					return OfValue[*option.Option[any]](option.Of(opt.GetValue()))
				} else {
					return OfValue[*option.Option[any]](option.None[any]())
				}
			}
		}
		return OfValue[*option.Option[any]](option.Of(this.GetValue()))

	}

	panic("Invalid empty result")
}

func (this *Result[T]) IsNil() bool {
	this.checkEvaluated()
	if this.IsOk() {
		return util.IsNil(this.Get())
	}
	return true
}

func (this *Result[T]) Failure() error {
	this.checkEvaluated()
	return this.failure.Get()
}

func (this *Result[T]) Get() T {
	return this.ok.Get()
}

func (this *Result[T]) OrNil() T {
	this.checkEvaluated()
	return this.ToOption().OrNil()
}

func (this *Result[T]) IfError(f func(error)) *Result[T] {
	this.checkEvaluated()
	if this.IsError() {
		f(this.failure.Get())
	}
	return this
}

func (this *Result[T]) IfOk(f func(T)) *Result[T] {
	this.checkEvaluated()
	if this.IsOk() {
		f(this.ok.Get())
	}
	return this
}

func (this *Result[T]) IfOkOpt(f func(*option.Option[T])) *Result[T] {
	this.checkEvaluated()
	if this.IsOk() {
		f(this.ToOption())
	}
	return this
}

func (this *Result[T]) IfOptEmpty(f func()) *Result[T] {
	this.checkEvaluated()
	if this.IsOk() && this.ToOption().IsEmpty() {
		f()
	}
	return this
}

func (this *Result[T]) IfOptNonEmpty(f func(T)) *Result[T] {
	this.checkEvaluated()
	if this.IsOk() && this.ToOption().NonEmpty() {
		f(this.Get())
	}
	return this
}

func (this *Result[T]) IsError() bool {
	this.checkEvaluated()
	return this.failure != nil
}

func (this *Result[T]) IsOk() bool {
	this.checkEvaluated()
	return this.ok != nil
}

func (this *Result[T]) IsEmpty() bool {
	this.checkEvaluated()
	return !this.IsError() && !this.IsOk()
}

func (this *Result[T]) IsResult() bool {
	return true
}

func (this *Result[T]) GetValue() interface{} {
	return this.Get()
}

func (this *Result[T]) GetError() error {
	return this.Failure()
}

func (this *Result[T]) Error() string {
	return this.Failure().Error()
}

func (this *Result[T]) HasError() bool {
	return this.IsError()
}

// FailWith if result is Ok and f() != nil, return new Result[T] with Failure(f())
func (this *Result[T]) FailWith(f func(T) error) *Result[T] {
	if this.IsOk() {
		if err := f(this.Get()); err != nil {
			return OfError[T](err)
		}
	}
	return this
}

func (this *Result[T]) Or(f func(T) T) *Result[T] {
	if this.IsOk() {
		return OfValue(f(this.Get()))
	}
	return this
}

func (this *Result[T]) OrElse(f func(T) *Result[T]) *Result[T] {
	if this.IsOk() {
		return f(this.Get())
	}
	return this
}

func (this *Result[T]) String() string {
	if this.IsError() {
		return fmt.Sprintf("Failure(%v)", this.Error())
	} else if this.IsOk() {
		return fmt.Sprintf("Ok(%v)", this.Get())
	} else {
		return fmt.Sprintf("Result(empty)")
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
	if v.IsOk() {
		if f(v.Get()) {
			return v
		}
	}
	return v
}

func Map[T any, R any](v *Result[T], f func(T) R) *Result[R] {
	if v.IsOk() {
		return OfValue(f(v.Get()))
	}
	return OfError[R](v.Failure())
}

func FlatMap[T any, R any](v *Result[T], f func(T) *Result[R]) *Result[R] {
	if v.IsOk() {
		return f(v.Get())
	}
	return OfError[R](v.Failure())
}
