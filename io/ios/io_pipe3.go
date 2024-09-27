package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/types/unit"
	"log"
	"reflect"
)

type IOPipe3[A, B, C, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C) *types.IO[T]
	fnResultOption func(A, B, C) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C) *result.Result[T]
	fnOption       func(A, B, C) *option.Option[T]
	fnValue        func(A, B, C) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe3IO[A, B, C, T any](f func(A, B, C) *types.IO[T]) *IOPipe3[A, B, C, T] {
	return &IOPipe3[A, B, C, T]{f: f}
}

func NewPipe3[A, B, C, T any](f func(A, B, C) *result.Result[*option.Option[T]]) *IOPipe3[A, B, C, T] {
	return &IOPipe3[A, B, C, T]{fnResultOption: f}
}

func NewPipe3OfValue[A, B, C, T any](f func(A, B, C) T) *IOPipe3[A, B, C, T] {
	return &IOPipe3[A, B, C, T]{fnValue: f}
}

func NewPipe3OfResult[A, B, C, T any](f func(A, B, C) *result.Result[T]) *IOPipe3[A, B, C, T] {
	return &IOPipe3[A, B, C, T]{fnResult: f}
}

func NewPipe3OfOption[A, B, C, T any](f func(A, B, C) *option.Option[T]) *IOPipe3[A, B, C, T] {
	return &IOPipe3[A, B, C, T]{fnOption: f}
}

func (this *IOPipe3[A, B, C, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe3[A, B, C, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOPipe3[A, B, C, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe3[A, B, C, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe3[A, B, C, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe3[A, B, C, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe3[A, B, C, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe3[A, B, C, T]) String() string {
	return fmt.Sprintf("Pipe3(%v)", this.value.String())
}

func (this *IOPipe3[A, B, C, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe3[A, B, C, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe3[A, B, C, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe3[A, B, C, T]) UnsafeRun() types.IOEffect {
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
		if this.f != nil {
			runnableIO := this.f(a, b, c)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b, c)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b, c))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b, c))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b, c)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
