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

type IOPipe6[A, B, C, D, E, F, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E, F) *types.IO[T]
	fnResultOption func(A, B, C, D, E, F) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E, F) *result.Result[T]
	fnOption       func(A, B, C, D, E, F) *option.Option[T]
	fnValue        func(A, B, C, D, E, F) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe6IO[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *types.IO[T]) *IOPipe6[A, B, C, D, E, F, T] {
	return &IOPipe6[A, B, C, D, E, F, T]{f: f}
}

func NewPipe6[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[*option.Option[T]]) *IOPipe6[A, B, C, D, E, F, T] {
	return &IOPipe6[A, B, C, D, E, F, T]{fnResultOption: f}
}

func NewPipe6OfValue[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) T) *IOPipe6[A, B, C, D, E, F, T] {
	return &IOPipe6[A, B, C, D, E, F, T]{fnValue: f}
}

func NewPipe6OfResult[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *result.Result[T]) *IOPipe6[A, B, C, D, E, F, T] {
	return &IOPipe6[A, B, C, D, E, F, T]{fnResult: f}
}

func NewPipe6OfOption[A, B, C, D, E, F, T any](f func(A, B, C, D, E, F) *option.Option[T]) *IOPipe6[A, B, C, D, E, F, T] {
	return &IOPipe6[A, B, C, D, E, F, T]{fnOption: f}
}

func (this *IOPipe6[A, B, C, D, E, F, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe6[A, B, C, D, E, F, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe6[A, B, C, D, E, F, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe6[A, B, C, D, E, F, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe6[A, B, C, D, E, F, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe6[A, B, C, D, E, F, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe6[A, B, C, D, E, F, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe6[A, B, C, D, E, F, T]) String() string {
	return fmt.Sprintf("Pipe6(%v)", this.value.String())
}

func (this *IOPipe6[A, B, C, D, E, F, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe6[A, B, C, D, E, F, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe6[A, B, C, D, E, F, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe6[A, B, C, D, E, F, T]) UnsafeRun() types.IOEffect {
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
		f := state.Consume[F](copyOfState)
		if this.f != nil {
			runnableIO := this.f(a, b, c, d, e, f)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c, d, e, f)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c, d, e, f))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c, d, e, f))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e, f)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
