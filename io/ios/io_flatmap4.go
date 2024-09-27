package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/types/unit"
	"reflect"
)

type IOFlatMap4[A, B, C, D, T any] struct {
	value      *result.Result[*option.Option[T]]
	prevEffect types.IOEffect
	f          func(A, B, C, D) *types.IO[T]
	ioA        *types.IO[A]
	ioB        *types.IO[B]
	ioC        *types.IO[C]
	ioD        *types.IO[D]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap4[A, B, C, D, T any](ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], ioD *types.IO[D], f func(A, B, C, D) *types.IO[T]) *IOFlatMap4[A, B, C, D, T] {
	return &IOFlatMap4[A, B, C, D, T]{f: f, ioA: ioA, ioB: ioB, ioC: ioC, ioD: ioD}
}

func (this *IOFlatMap4[A, B, C, D, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOFlatMap4[A, B, C, D, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap4[A, B, C, D, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOFlatMap4[A, B, C, D, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFlatMap4[A, B, C, D, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap4[A, B, C, D, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFlatMap4[A, B, C, D, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFlatMap4[A, B, C, D, T]) String() string {
	return fmt.Sprintf("FlatMap4(%v)", this.value.String())
}

func (this *IOFlatMap4[A, B, C, D, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap4[A, B, C, D, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap4[A, B, C, D, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap4[A, B, C, D, T]) UnsafeRun() types.IOEffect {
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
		runnableIO := NewFlatMap3[A, B, C, T](
			this.ioA, this.ioB, this.ioC, func(a A, b B, c C) *types.IO[T] {
				return NewFlatMapIO[D, T](
					this.ioD,
					func(d D) *types.IO[T] {
						return this.f(a, b, c, d)
					}).Lift()
			}).Lift()
		this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(types.IOEffect)
}
