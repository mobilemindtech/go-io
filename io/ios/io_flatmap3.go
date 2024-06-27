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

type IOFlatMap3[A, B, C, T any] struct {
	value      *result.Result[*option.Option[T]]
	prevEffect types.IOEffect
	f          func(A, B, C) *types.IO[T]
	ioA        *types.IO[A]
	ioB        *types.IO[B]
	ioC        *types.IO[C]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap3[A, B, C, T any](ioA *types.IO[A], ioB *types.IO[B], ioC *types.IO[C], f func(A, B, C) *types.IO[T]) *IOFlatMap3[A, B, C, T] {
	return &IOFlatMap3[A, B, C, T]{f: f, ioA: ioA, ioB: ioB, ioC: ioC}
}

func (this *IOFlatMap3[A, B, C, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOFlatMap3[A, B, C, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap3[A, B, C, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOFlatMap3[A, B, C, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFlatMap3[A, B, C, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap3[A, B, C, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFlatMap3[A, B, C, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFlatMap3[A, B, C, T]) String() string {
	return fmt.Sprintf("FlatMap3(%v)", this.value.String())
}

func (this *IOFlatMap3[A, B, C, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap3[A, B, C, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap3[A, B, C, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap3[A, B, C, T]) UnsafeRun() types.IOEffect {
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
		runnableIO := NewFlatMap2[A, B, T](
			this.ioA, this.ioB, func(a A, b B) *types.IO[T] {
				return NewFlatMapIO[C, T](
					this.ioC, func(c C) *types.IO[T] {
						return this.f(a, b, c)
					}).Lift()
			}).Lift()
		this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(types.IOEffect)
}
