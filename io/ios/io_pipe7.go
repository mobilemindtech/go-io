package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/runtime"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"log"
	"reflect"
)

type IOPipe7[A, B, C, D, E, F, G, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E, F, G) *types.IO[T]
	fnResultOption func(A, B, C, D, E, F, G) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E, F, G) *result.Result[T]
	fnOption       func(A, B, C, D, E, F, G) *option.Option[T]
	fnValue        func(A, B, C, D, E, F, G) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe7IO[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *types.IO[T]) *IOPipe7[A, B, C, D, E, F, G, T] {
	return &IOPipe7[A, B, C, D, E, F, G, T]{f: f}
}

func NewPipe7[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[*option.Option[T]]) *IOPipe7[A, B, C, D, E, F, G, T] {
	return &IOPipe7[A, B, C, D, E, F, G, T]{fnResultOption: f}
}

func NewPipe7OfValue[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) T) *IOPipe7[A, B, C, D, E, F, G, T] {
	return &IOPipe7[A, B, C, D, E, F, G, T]{fnValue: f}
}

func NewPipe7OfResult[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *result.Result[T]) *IOPipe7[A, B, C, D, E, F, G, T] {
	return &IOPipe7[A, B, C, D, E, F, G, T]{fnResult: f}
}

func NewPipe7OfOption[A, B, C, D, E, F, G, T any](f func(A, B, C, D, E, F, G) *option.Option[T]) *IOPipe7[A, B, C, D, E, F, G, T] {
	return &IOPipe7[A, B, C, D, E, F, G, T]{fnOption: f}
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) String() string {
	return fmt.Sprintf("Pipe7(%v)", this.value.String())
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe7[A, B, C, D, E, F, G, T]) UnsafeRun() types.IOEffect {
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
		g := state.Consume[G](copyOfState)
		if this.f != nil {
			runnableIO := this.f(a, b, c, d, e, f, g)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c, d, e, f, g)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c, d, e, f, g))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c, d, e, f, g))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e, f, g)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
