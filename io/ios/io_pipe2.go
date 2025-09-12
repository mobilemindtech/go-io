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

type IOPipe2[A, B, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B) *types.IO[T]
	fnResultOption func(A, B) *result.Result[*option.Option[T]]
	fnResult       func(A, B) *result.Result[T]
	fnOption       func(A, B) *option.Option[T]
	fnValue        func(A, B) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe2IO[A, B, T any](f func(A, B) *types.IO[T]) *IOPipe2[A, B, T] {
	return &IOPipe2[A, B, T]{f: f}
}

func NewPipe2[A, B, T any](f func(A, B) *result.Result[*option.Option[T]]) *IOPipe2[A, B, T] {
	return &IOPipe2[A, B, T]{fnResultOption: f}
}

func NewPipe2OfValue[A, B, T any](f func(A, B) T) *IOPipe2[A, B, T] {
	return &IOPipe2[A, B, T]{fnValue: f}
}

func NewPipe2OfResult[A, B, T any](f func(A, B) *result.Result[T]) *IOPipe2[A, B, T] {
	return &IOPipe2[A, B, T]{fnResult: f}
}

func NewPipe2OfOption[A, B, T any](f func(A, B) *option.Option[T]) *IOPipe2[A, B, T] {
	return &IOPipe2[A, B, T]{fnOption: f}
}

func (this *IOPipe2[A, B, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe2[A, B, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOPipe2[A, B, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe2[A, B, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe2[A, B, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe2[A, B, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe2[A, B, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe2[A, B, T]) String() string {
	return fmt.Sprintf("Pipe2(%v)", this.value.String())
}

func (this *IOPipe2[A, B, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe2[A, B, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe2[A, B, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe2[A, B, T]) UnsafeRun() types.IOEffect {
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
		if this.f != nil {
			runnableIO := this.f(a, b)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a, b)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a, b))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a, b))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a, b)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
