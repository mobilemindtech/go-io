package types

import (
	"github.com/mobilemindtec/go-io/collections"
	"github.com/mobilemindtec/go-io/result"
)

type IIO interface {
	UnsafeRunIO() *result.Result[any]
}

type IO[T any] struct {
	stack   *collections.Stack[IOEffect]
	varName string
}

func NewIO[T any]() *IO[T] {
	return &IO[T]{stack: collections.NewStack[IOEffect]()}
}

func (this *IO[T]) push(val IOEffect) *IO[T] {
	this.stack.
		Peek().
		IfNonEmpty(func(eff IOEffect) {
			val.SetPrevEffect(eff)
		})
	this.stack.Push(val)
	return this
}

func (this *IO[T]) As(name string) *IO[T] {
	this.varName = name
	return this
}

func (this *IO[T]) Effect(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Effects(vals ...IOEffect) *IO[T] {
	for _, eff := range vals {
		this.Effects(eff)
	}
	return this
}

func (this *IO[T]) Pure(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Map(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) FlatMap(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Recover(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) FailWith(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Filter(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Tap(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Debug(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) SliceForeach(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) SliceMap(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) SliceFlatMap(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) SliceFilter(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) Attempt(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func runStackIO[T any](stack *collections.Stack[IOEffect], currEff IOEffect) IOEffect {
	if stack.IsNonEmpty() {
		eff := stack.UnsafePop()
		runStackIO[T](stack, eff)
	}
	return currEff.UnsafeRun()
}

func (this *IO[T]) UnsafeRun() *result.Result[T] {
	var effResult IOEffect
	if this.stack.Count() > 1 {
		effResult = runStackIO[T](this.stack, this.stack.UnsafePop())
	} else {
		effResult = this.stack.UnsafePop().UnsafeRun()
	}
	r := effResult.GetResult()
	return result.Make[T](r.OrNil().(T), r.Error())
}

func (this *IO[T]) UnsafeRunIO() *result.Result[any] {
	return this.UnsafeRun().ToResultOfAny()
}

func (this *IO[T]) GetVarName() string {
	return this.varName
}
