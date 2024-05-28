package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/collections"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IO[T any] struct {
	stack   *collections.Stack[IOEffect]
	varName string
	state   *state.State
	debug   bool
}

func NewIO[T any]() *IO[T] {
	return &IO[T]{stack: collections.NewStack[IOEffect](), state: state.NewState()}
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
		this.Effect(eff)
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

func (this *IO[T]) MaybeFail(val IOEffect) *IO[T] {
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

func (this *IO[T]) Or(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) OrElse(val IOEffect) *IO[T] {
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

func (this *IO[T]) Exec(val IOEffect) *IO[T] {
	this.push(val)
	return this
}

func (this *IO[T]) runStackIO(currEff IOEffect) IOEffect {
	if this.stack.IsNonEmpty() {
		eff := this.stack.UnsafePop()
		this.runStackIO(eff)
	}

	if this.debug {
		log.Printf("IO>> UnsafeRun %v", reflect.TypeOf(currEff).Name())
	}

	currEff.SetState(this.state)
	currEff.SetDebug(this.debug)

	r := currEff.UnsafeRun()

	if this.debug {
		log.Printf("IO>> UnsafeRun %v = %v", reflect.TypeOf(currEff).Name(), r.String())
	}

	return r
}

func (this *IO[T]) UnsafeRun() *result.Result[*option.Option[T]] {
	effResult := this.runStackIO(this.stack.UnsafePop())
	r := effResult.GetResult()

	if r.IsError() {
		return result.OfError[*option.Option[T]](r.GetError())
	}

	if util.CanNil(reflect.ValueOf(r.GetValue()).Kind()) && util.IsNil(r.GetValue()) || r.Get().Empty() {
		return result.OfValue(option.None[T]())
	}

	val := r.Get().GetValue()
	if v, ok := val.(T); ok {
		return result.OfValue(option.Some(v))
	}
	typOf := reflect.TypeFor[T]()
	panic(fmt.Sprintf("can't cast %v to IO result type %v", r.GetValue(), typOf))
}

func (this *IO[T]) SetState(st *state.State) {
	this.state = st
}

func (this *IO[T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IO[T]) DebugOn() *IO[T] {
	this.SetDebug(true)
	return this
}

func (this *IO[T]) UnsafeRunIO() ResultOptionAny {
	return this.UnsafeRun().ToResultOfOption()
}

func (this *IO[T]) GetVarName() string {
	return this.varName
}
