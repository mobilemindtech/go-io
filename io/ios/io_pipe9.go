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

type IOPipe9[A, B, C, D, E, F, G, H, I, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E, F, G, H, I) *types.IO[T]
	fnResultOption func(A, B, C, D, E, F, G, H, I) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E, F, G, H, I) *result.Result[T]
	fnOption       func(A, B, C, D, E, F, G, H, I) *option.Option[T]
	fnValue        func(A, B, C, D, E, F, G, H, I) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe9IO[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *types.IO[T]) *IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return &IOPipe9[A, B, C, D, E, F, G, H, I, T]{f: f}
}

func NewPipe9[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[*option.Option[T]]) *IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return &IOPipe9[A, B, C, D, E, F, G, H, I, T]{fnResultOption: f}
}

func NewPipe9OfValue[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) T) *IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return &IOPipe9[A, B, C, D, E, F, G, H, I, T]{fnValue: f}
}

func NewPipe9OfResult[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *result.Result[T]) *IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return &IOPipe9[A, B, C, D, E, F, G, H, I, T]{fnResult: f}
}

func NewPipe9OfOption[A, B, C, D, E, F, G, H, I, T any](f func(A, B, C, D, E, F, G, H, I) *option.Option[T]) *IOPipe9[A, B, C, D, E, F, G, H, I, T] {
	return &IOPipe9[A, B, C, D, E, F, G, H, I, T]{fnOption: f}
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) String() string {
	return fmt.Sprintf("Pipe9(%v)", this.value.String())
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe9[A, B, C, D, E, F, G, H, I, T]) UnsafeRun() types.IOEffect {
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
		i := state.Consume[I](copyOfState)
		if this.f != nil {
			runnableIO := this.f(a, b, c, d, e, f, g, h, i)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c, d, e, f, g, h, i)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c, d, e, f, g, h, i))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c, d, e, f, g, h, i))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e, f, g, h, i)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
