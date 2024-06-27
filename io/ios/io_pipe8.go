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

type IOPipe8[A, B, C, D, E, F, G, H, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E, F, G, H) *types.IO[T]
	fnResultOption func(A, B, C, D, E, F, G, H) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E, F, G, H) *result.Result[T]
	fnOption       func(A, B, C, D, E, F, G, H) *option.Option[T]
	fnValue        func(A, B, C, D, E, F, G, H) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe8IO[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *types.IO[T]) *IOPipe8[A, B, C, D, E, F, G, H, T] {
	return &IOPipe8[A, B, C, D, E, F, G, H, T]{f: f}
}

func NewPipe8[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[*option.Option[T]]) *IOPipe8[A, B, C, D, E, F, G, H, T] {
	return &IOPipe8[A, B, C, D, E, F, G, H, T]{fnResultOption: f}
}

func NewPipe8OfValue[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) T) *IOPipe8[A, B, C, D, E, F, G, H, T] {
	return &IOPipe8[A, B, C, D, E, F, G, H, T]{fnValue: f}
}

func NewPipe8OfResult[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *result.Result[T]) *IOPipe8[A, B, C, D, E, F, G, H, T] {
	return &IOPipe8[A, B, C, D, E, F, G, H, T]{fnResult: f}
}

func NewPipe8OfOption[A, B, C, D, E, F, G, H, T any](f func(A, B, C, D, E, F, G, H) *option.Option[T]) *IOPipe8[A, B, C, D, E, F, G, H, T] {
	return &IOPipe8[A, B, C, D, E, F, G, H, T]{fnOption: f}
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) String() string {
	return fmt.Sprintf("Pipe8(%v)", this.value.String())
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe8[A, B, C, D, E, F, G, H, T]) UnsafeRun() types.IOEffect {
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
		h := state.Consume[H](copyOfState)
		if this.f != nil {
			runnableIO := this.f(a, b, c, d, e, f, g, h)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c, d, e, f, g, h)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c, d, e, f, g, h))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c, d, e, f, g, h))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e, f, g, h)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
