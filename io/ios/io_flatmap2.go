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

type IOFlatMap2[A, B, T any] struct {
	value      *result.Result[*option.Option[T]]
	prevEffect types.IOEffect
	f          func(A, B) *types.IO[T]
	ioA        *types.IO[A]
	ioB        *types.IO[B]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap2[A, B, T any](ioA *types.IO[A], ioB *types.IO[B], f func(A, B) *types.IO[T]) *IOFlatMap2[A, B, T] {
	return &IOFlatMap2[A, B, T]{f: f, ioA: ioA, ioB: ioB}
}

func (this *IOFlatMap2[A, B, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOFlatMap2[A, B, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap2[A, B, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOFlatMap2[A, B, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFlatMap2[A, B, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap2[A, B, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFlatMap2[A, B, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFlatMap2[A, B, T]) String() string {
	return fmt.Sprintf("FlatMap2(%v)", this.value.String())
}

func (this *IOFlatMap2[A, B, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap2[A, B, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap2[A, B, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap2[A, B, T]) UnsafeRun() types.IOEffect {
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
		runnableIO := NewFlatMapIO[A, T](this.ioA, func(a A) *types.IO[T] {
			return NewFlatMapIO[B, T](this.ioB, func(b B) *types.IO[T] {
				return this.f(a, b)
			}).Lift()
		}).Lift()
		this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(types.IOEffect)
}
