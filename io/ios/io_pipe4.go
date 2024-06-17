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

type IOPipe4[A, B, C, D, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D) types.IORunnable
	fnResultOption func(A, B, C, D) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D) *result.Result[T]
	fnOption       func(A, B, C, D) *option.Option[T]
	fnValue        func(A, B, C, D) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe4IO[A, B, C, D, T any](f func(A, B, C, D) types.IORunnable) *IOPipe4[A, B, C, D, T] {
	return &IOPipe4[A, B, C, D, T]{f: f}
}

func NewPipe4[A, B, C, D, T any](f func(A, B, C, D) *result.Result[*option.Option[T]]) *IOPipe4[A, B, C, D, T] {
	return &IOPipe4[A, B, C, D, T]{fnResultOption: f}
}

func NewPipe4OfValue[A, B, C, D, T any](f func(A, B, C, D) T) *IOPipe4[A, B, C, D, T] {
	return &IOPipe4[A, B, C, D, T]{fnValue: f}
}

func NewPipe4OfResult[A, B, C, D, T any](f func(A, B, C, D) *result.Result[T]) *IOPipe4[A, B, C, D, T] {
	return &IOPipe4[A, B, C, D, T]{fnResult: f}
}

func NewPipe4OfOption[A, B, C, D, T any](f func(A, B, C, D) *option.Option[T]) *IOPipe4[A, B, C, D, T] {
	return &IOPipe4[A, B, C, D, T]{fnOption: f}
}

func (this *IOPipe4[A, B, C, D, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe4[A, B, C, D, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe4[A, B, C, D, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe4[A, B, C, D, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe4[A, B, C, D, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe4[A, B, C, D, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe4[A, B, C, D, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe4[A, B, C, D, T]) String() string {
	return fmt.Sprintf("Pipe4(%v)", this.value.String())
}

func (this *IOPipe4[A, B, C, D, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe4[A, B, C, D, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe4[A, B, C, D, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe4[A, B, C, D, T]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[T]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[T]](r.Failure())
		} else {
			copyOfState := this.state.Copy()
			a := state.Consume[A](copyOfState)
			b := state.Consume[B](copyOfState)
			c := state.Consume[C](copyOfState)
			d := state.Consume[D](copyOfState)
			if this.f != nil {
				runnableIO := this.f(a, b, c, d)
				this.value = runtime.New[T](runnableIO).UnsafeRun()
			} else if this.fnResultOption != nil {
				this.value = this.fnResultOption(a, b, c, d)
			} else if this.fnOption != nil {
				this.value = result.OfValue(this.fnOption(a, b, c, d))
			} else if this.fnResult != nil {
				this.value = ResultToResultOption(this.fnResult(a, b, c, d))
			} else if this.fnValue != nil {
				this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d)))
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
