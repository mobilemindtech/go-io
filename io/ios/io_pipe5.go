package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOPipe5[A, B, C, D, E, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E) *types.IO[T]
	fnResultOption func(A, B, C, D, E) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E) *result.Result[T]
	fnOption       func(A, B, C, D, E) *option.Option[T]
	fnValue        func(A, B, C, D, E) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe5IO[A, B, C, D, E, T any](f func(A, B, C, D, E) *types.IO[T]) *IOPipe5[A, B, C, D, E, T] {
	return &IOPipe5[A, B, C, D, E, T]{f: f}
}

func NewPipe5[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[*option.Option[T]]) *IOPipe5[A, B, C, D, E, T] {
	return &IOPipe5[A, B, C, D, E, T]{fnResultOption: f}
}

func NewPipe5OfValue[A, B, C, D, E, T any](f func(A, B, C, D, E) T) *IOPipe5[A, B, C, D, E, T] {
	return &IOPipe5[A, B, C, D, E, T]{fnValue: f}
}

func NewPipe5OfResult[A, B, C, D, E, T any](f func(A, B, C, D, E) *result.Result[T]) *IOPipe5[A, B, C, D, E, T] {
	return &IOPipe5[A, B, C, D, E, T]{fnResult: f}
}

func NewPipe5OfOption[A, B, C, D, E, T any](f func(A, B, C, D, E) *option.Option[T]) *IOPipe5[A, B, C, D, E, T] {
	return &IOPipe5[A, B, C, D, E, T]{fnOption: f}
}

func (this *IOPipe5[A, B, C, D, E, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe5[A, B, C, D, E, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe5[A, B, C, D, E, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe5[A, B, C, D, E, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe5[A, B, C, D, E, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe5[A, B, C, D, E, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe5[A, B, C, D, E, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe5[A, B, C, D, E, T]) String() string {
	return fmt.Sprintf("Pipe5(%v)", this.value.String())
}

func (this *IOPipe5[A, B, C, D, E, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe5[A, B, C, D, E, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe5[A, B, C, D, E, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe5[A, B, C, D, E, T]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[T]())
	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[T]](r.Failure())
			execute = false
		}
	}

	if execute {
		copyOfState := this.state.Copy()
		a := state.Consume[A](copyOfState)
		b := state.Consume[B](copyOfState)
		c := state.Consume[C](copyOfState)
		d := state.Consume[D](copyOfState)
		e := state.Consume[E](copyOfState)
		if this.f != nil {
			runnableIO := this.f(a, b, c, d, e)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c, d, e)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c, d, e))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c, d, e))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
