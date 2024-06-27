package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"reflect"
)

type IOFlatMap5[A, B, C, D, E, T any] struct {
	value      *result.Result[*option.Option[T]]
	prevEffect types.IOEffect
	f          func(A, B, C, D, E) *types.IO[T]
	ioA        *types.IO[A]
	ioB        *types.IO[B]
	ioC        *types.IO[C]
	ioD        *types.IO[D]
	ioE        *types.IO[E]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap5[A, B, C, D, E, T any](
	ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], ioD *types.IO[D], ioE *types.IO[E], f func(A, B, C, D, E) *types.IO[T]) *IOFlatMap5[A, B, C, D, E, T] {
	return &IOFlatMap5[A, B, C, D, E, T]{f: f, ioA: ioA, ioB: ioB, ioC: ioC, ioD: ioD, ioE: ioE}
}

func (this *IOFlatMap5[A, B, C, D, E, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOFlatMap5[A, B, C, D, E, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap5[A, B, C, D, E, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOFlatMap5[A, B, C, D, E, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFlatMap5[A, B, C, D, E, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap5[A, B, C, D, E, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFlatMap5[A, B, C, D, E, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFlatMap5[A, B, C, D, E, T]) String() string {
	return fmt.Sprintf("FlatMap5(%v)", this.value.String())
}

func (this *IOFlatMap5[A, B, C, D, E, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap5[A, B, C, D, E, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap5[A, B, C, D, E, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap5[A, B, C, D, E, T]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[T]())
	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[T]](r.Failure())
			execute = false
		} else if r.Get().Empty() {
			execute = false
		}
	}

	if execute {
		runnableIO := NewFlatMap4[A, B, C, D, T](
			this.ioA, this.ioB, this.ioC, this.ioD, func(a A, b B, c C, d D) *types.IO[T] {
				return NewFlatMapIO[E, T](
					this.ioE,
					func(e E) *types.IO[T] {
						return this.f(a, b, c, d, e)
					}).Lift()
			}).Lift()
		this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(types.IOEffect)
}
