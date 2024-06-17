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

type IOPipe10[A, B, C, D, E, F, G, H, I, J, T any] struct {
	value          *result.Result[*option.Option[T]]
	prevEffect     types.IOEffect
	f              func(A, B, C, D, E, F, G, H, I, J) types.IORunnable
	fnResultOption func(A, B, C, D, E, F, G, H, I, J) *result.Result[*option.Option[T]]
	fnResult       func(A, B, C, D, E, F, G, H, I, J) *result.Result[T]
	fnOption       func(A, B, C, D, E, F, G, H, I, J) *option.Option[T]
	fnValue        func(A, B, C, D, E, F, G, H, I, J) T
	state          *state.State
	debug          bool
	debugInfo      *types.IODebugInfo
}

func NewPipe10IO[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) types.IORunnable) *IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return &IOPipe10[A, B, C, D, E, F, G, H, I, J, T]{f: f}
}

func NewPipe10[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[*option.Option[T]]) *IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return &IOPipe10[A, B, C, D, E, F, G, H, I, J, T]{fnResultOption: f}
}

func NewPipe10OfValue[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) T) *IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return &IOPipe10[A, B, C, D, E, F, G, H, I, J, T]{fnValue: f}
}

func NewPipe10OfResult[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *result.Result[T]) *IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return &IOPipe10[A, B, C, D, E, F, G, H, I, J, T]{fnResult: f}
}

func NewPipe10OfOption[A, B, C, D, E, F, G, H, I, J, T any](f func(A, B, C, D, E, F, G, H, I, J) *option.Option[T]) *IOPipe10[A, B, C, D, E, F, G, H, I, J, T] {
	return &IOPipe10[A, B, C, D, E, F, G, H, I, J, T]{fnOption: f}
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) TypeIn() reflect.Type {
	return reflect.TypeFor[*types.Unit]()
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) SetState(st *state.State) {
	this.state = st
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) String() string {
	return fmt.Sprintf("Pipe10(%v)", this.value.String())
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOPipe10[A, B, C, D, E, F, G, H, I, J, T]) UnsafeRun() types.IOEffect {
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
			e := state.Consume[E](copyOfState)
			f := state.Consume[F](copyOfState)
			g := state.Consume[G](copyOfState)
			h := state.Consume[H](copyOfState)
			i := state.Consume[I](copyOfState)
			j := state.Consume[J](copyOfState)
			if this.f != nil {
				runnableIO := this.f(a, b, c, d, e, f, g, h, i, j)
				this.value = runtime.New[T](runnableIO).UnsafeRun()
			} else if this.fnResultOption != nil {
				this.value = this.fnResultOption(a, b, c, d, e, f, g, h, i, j)
			} else if this.fnOption != nil {
				this.value = result.OfValue(this.fnOption(a, b, c, d, e, f, g, h, i, j))
			} else if this.fnResult != nil {
				this.value = ResultToResultOption(this.fnResult(a, b, c, d, e, f, g, h, i, j))
			} else if this.fnValue != nil {
				this.value = result.OfValue(option.Of(this.fnValue(a, b, c, d, e, f, g, h, i, j)))
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
