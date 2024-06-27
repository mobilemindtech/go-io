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

type IOPipe[A, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A) *types.IO[T]
	fnResultOption func(A) *result.Result[*option.Option[T]]
	fnResult       func(A) *result.Result[T]
	fnOption       func(A) *option.Option[T]
	fnValue        func(A) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipeIO[A, T any](f func(A) *types.IO[T]) *IOPipe[A, T] {
	return &IOPipe[A, T]{f: f}
}

func NewPipe[A, T any](f func(A) *result.Result[*option.Option[T]]) *IOPipe[A, T] {
	return &IOPipe[A, T]{fnResultOption: f}
}

func NewPipeOfValue[A, T any](f func(A) T) *IOPipe[A, T] {
	return &IOPipe[A, T]{fnValue: f}
}

func NewPipeOfResult[A, T any](f func(A) *result.Result[T]) *IOPipe[A, T] {
	return &IOPipe[A, T]{fnResult: f}
}

func NewPipeOfOption[A, T any](f func(A) *option.Option[T]) *IOPipe[A, T] {
	return &IOPipe[A, T]{fnOption: f}
}

func (this *IOPipe[A, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe[A, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe[A, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe[A, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe[A, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe[A, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe[A, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe[A, T]) String() string {
	return fmt.Sprintf("Pipe(%v)", this.value.String())
}

func (this *IOPipe[A, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe[A, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe[A, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe[A, T]) UnsafeRun() types.IOEffect {
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
		a := state.Var[A](this.state)

		if this.f != nil {
			runnableIO := this.f(a)
			this.value = runtime.NewWithState[T](this.state, runnableIO).UnsafeRun()
		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption(a)
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption(a))
		} else if this.fnResult != nil {
			this.value = ResultToResultOption(this.fnResult(a))
		} else if this.fnValue != nil {
			this.value = result.OfValue(option.Of(this.fnValue(a)))
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
