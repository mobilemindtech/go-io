package result

import (
	"fmt"
	"reflect"

	"github.com/mobilemindtech/go-io/fault"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/types/unit"
	"github.com/mobilemindtech/go-io/util"
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

func Try[T any](f func() (T, error)) (ret *Result[T]) {

	defer func() {
		if err := recover(); err != nil {
			ret = OfError[T](fault.AnyToError(err))
		}
	}()

	v, e := f()
	return Make(v, e)
}

func TryUnit(f func() error) (ret *Result[*unit.Unit]) {
	defer func() {
		if err := recover(); err != nil {
			ret = OfError[*unit.Unit](fault.AnyToError(err))
		}
	}()

	e := f()
	return Make(unit.OfUnit(), e)
}

func TryMap[A, B any](ftry func() (A, error), f func(A) B) *Result[B] {
	v, e := ftry()
	res := Make(v, e)

	if res.IsError() {
		return OfError[B](res.Failure())
	}
	return OfValue(f(res.Get()))

}

func TryFlatMap[A, B any](ftry func() (A, error), f func(A) *Result[B]) *Result[B] {
	v, e := ftry()
	res := Make(v, e)

	if res.IsError() {
		return OfError[B](res.Failure())
	}
	return f(res.Get())

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

func OfNone[T any]() *Result[*option.Option[T]] {
	return &Result[*option.Option[T]]{ok: _newOk(option.None[T]()), evaluated: true}
}

func OfSome[T any](val T) *Result[*option.Option[T]] {
	return &Result[*option.Option[T]]{ok: _newOk(option.Of(val)), evaluated: true}
}

func OfErrorOption[T any](err error) *Result[*option.Option[T]] {
	return &Result[*option.Option[T]]{failure: _newFailure(err), evaluated: true}
}

func MapToResultOption[T any](res *Result[T]) *Result[*option.Option[T]] {

	if res.IsError() {
		return OfErrorOption[T](res.Failure())
	}
	return &Result[*option.Option[T]]{ok: _newOk(option.Of(res.Get())), evaluated: true}
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

func OfErrorf[T any](msg string, args ...any) *Result[T] {
	return OfError[T](fmt.Errorf(msg, args...))
}

func OfErrorOrValue[T any](err error, def T) *Result[T] {
	if err != nil {
		return OfError[T](err)
	}
	return OfValue(def)
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

func (this *Result[T]) FailureOrNil() error {
	this.checkEvaluated()
	if this.HasError() {
		return this.Failure()
	}
	return nil
}

func (this *Result[T]) Unsafe() T {
	if this.IsError() {
		panic(this.Failure())
	}
	return this.Get()
}

func (this *Result[T]) Get() T {
	return this.ok.Get()
}

func (this *Result[T]) OrNil() T {
	this.checkEvaluated()
	return this.ToOption().OrNil()
}

func (this *Result[T]) OrPanic(msg string) T {
	this.checkEvaluated()

	if this.IsOk() {
		return this.Get()
	}

	panic(fmt.Sprintf(msg, this.GetError()))
}

func (this *Result[T]) IfError(f func(error)) *Result[T] {
	this.checkEvaluated()
	if this.IsError() {
		f(this.failure.Get())
	}
	return this
}

func (this *Result[T]) Resolve(ferr func(error), fok func(T)) *Result[T] {
	this.checkEvaluated()
	if this.IsError() {
		ferr(this.failure.Get())
	} else {
		fok(this.Get())
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

func (this *Result[T]) Foreach(f func(T)) *Result[T] {
	if this.IsOk() {
		f(this.Get())
	}
	return this
}

func (this *Result[T]) Update(f func(T) T) *Result[T] {
	if this.IsOk() {
		return OfValue(f(this.Get()))
	}
	return this
}

func (this *Result[T]) Exec(f func(T) *Result[T]) *Result[T] {
	if this.IsOk() {
		return f(this.Get())
	}
	return this
}

func (this *Result[T]) TryExec(f func(T) error) *Result[T] {
	if this.IsOk() {

		if err := f(this.Get()); err != nil {
			return OfError[T](err)
		}

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

func (this *Result[T]) ErrorComplement(f func(error) error) *Result[T] {
	if this.IsError() {
		return OfError[T](f(this.Failure()))
	}
	return this
}

func (this *Result[T]) CatchAll(f func(error) *Result[T]) *Result[T] {
	if this.IsError() {
		return f(this.Failure())
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

func (this *Result[T]) MapToBool() *Result[bool] {
	if this.IsError() {
		return OfError[bool](this.Failure())
	}
	return OfValue(true)
}

func (this *Result[T]) MapToUnit() *Result[*unit.Unit] {
	if this.IsError() {
		return OfError[*unit.Unit](this.Failure())
	}
	return OfValue(unit.OfUnit())
}

func (this *Result[T]) MapToBoolWith(f func(T) *Result[bool]) *Result[bool] {
	if this.IsError() {
		return OfError[bool](this.Failure())
	}
	return f(this.Get())
}

func (this *Result[T]) MapToUnitWith(f func(T) *Result[*unit.Unit]) *Result[*unit.Unit] {
	if this.IsError() {
		return OfError[*unit.Unit](this.Failure())
	}
	return f(this.Get())
}

func (this *Result[T]) ErrorOrNil() error {
	if this.IsError() {
		return this.GetError()
	}
	return nil
}

func (this *Result[T]) RaiseWhen(err error, f func(T) bool) *Result[T] {
	if !this.IsError() && f(this.Get()) {
		return OfError[T](err)
	}
	return this
}

func (this *Result[T]) UnwrapTo(f func(interface{})) *Result[T] {
	if this.IsError() {
		f(this.Failure())
	} else {
		f(this.Get())
	}
	return this
}

func (this *Result[T]) PanicIfFail() *Result[T] {
	if this.IsError() {
		panic(this.Failure())
	}
	return this
}

func (this *Result[T]) FilterOrError(f func(T) bool, err error) *Result[T] {
	if this.IsOk() {
		if f(this.Get()) {
			return this
		} else {
			return OfError[T](err)
		}
	}
	return this
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

func Map[T, R any](v *Result[T], f func(T) R) *Result[R] {
	if v.IsOk() {
		return OfValue(f(v.Get()))
	}
	return OfError[R](v.Failure())
}

func FlatMapOption[T, R any](v *Result[*option.Option[T]], f func(T) *option.Option[R]) *Result[*option.Option[R]] {
	if v.IsError() {
		return OfError[*option.Option[R]](v.Failure())
	}
	if v.IsOk() && v.Get().IsSome() {
		return OfValue(f(v.Get().Get()))
	}
	return OfValue(option.None[R]())

}

func MapOptionToValue[T, R any](v *Result[*option.Option[T]], f func(T) R, orElseVal R) *Result[R] {
	if v.IsError() {
		return OfError[R](v.Failure())
	}
	if v.IsOk() && v.Get().IsSome() {
		r := f(v.Get().Get())
		return OfValue(r)
	}
	return OfValue(orElseVal)
}

func UnwapOptionValueOrNil[T any](v *Result[*option.Option[T]]) T {
	if v.IsOk() && v.Get().IsSome() {
		return v.Get().Get()
	}
	var t T
	return t
}

func FlatMap[T, R any](v *Result[T], f func(T) *Result[R]) *Result[R] {
	if v.IsOk() {
		return f(v.Get())
	}
	return OfError[R](v.Failure())
}

func MapToValue[A, B any](res *Result[A], b B) *Result[B] {
	if res.HasError() {
		return OfError[B](res.Failure())
	}
	return OfValue(b)
}

func MapToValueOfOption[A, B any](res *Result[*option.Option[A]], b B) *Result[*option.Option[B]] {
	if res.HasError() {
		return OfError[*option.Option[B]](res.Failure())
	}
	if res.Get().IsNone() {
		return OfValue(option.None[B]())
	}
	return OfValue(option.Some(b))
}

func SliceFlatMap[A, B any](vs []A, f func(A) *Result[B]) *Result[[]B] {
	var items []B

	for _, it := range vs {
		res := f(it)
		if res.IsOk() {
			items = append(items, res.Get())
		} else {
			return OfError[[]B](res.GetError())
		}
	}

	return OfValue(items)
}

func All(results ...IResult) *Result[*unit.Unit] {
	for _, it := range results {
		if it.HasError() {
			return OfError[*unit.Unit](it.GetError())
		}
	}
	return OfValue(unit.OfUnit())
}

func AllReturnFirst[T any](results ...*Result[T]) *Result[T] {
	for _, it := range results {
		if it.HasError() {
			return OfError[T](it.GetError())
		}
	}
	return results[0]
}

func AllReturnLast[T any](results ...*Result[T]) *Result[T] {
	for _, it := range results {
		if it.HasError() {
			return OfError[T](it.GetError())
		}
	}
	return results[len(results)-1]
}
