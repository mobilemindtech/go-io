package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/runtime"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/types/unit"
	"github.com/mobilemindtech/go-io/util"
	"reflect"
)

type IOFlatMap[A any, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect types.IOEffect
	f          func(A) *types.IO[B]
	ioA        *types.IO[A]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap[A any, B any](f func(A) *types.IO[B]) *IOFlatMap[A, B] {
	return &IOFlatMap[A, B]{f: f}
}

func NewFlatMapIO[A any, B any](ioA *types.IO[A], f func(A) *types.IO[B]) *IOFlatMap[A, B] {
	return &IOFlatMap[A, B]{f: f, ioA: ioA}
}

func (this *IOFlatMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap[A, B]) TypeIn() reflect.Type {

	if this.ioA != nil {
		return reflect.TypeFor[*unit.Unit]()
	} else {
		return reflect.TypeFor[A]()
	}
}

func (this *IOFlatMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[B]()
}

func (this *IOFlatMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFlatMap[A, B]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFlatMap[A, B]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFlatMap[A, B]) String() string {
	return fmt.Sprintf("FlatMap(%v)", this.value.String())
}

func (this *IOFlatMap[A, B]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFlatMap[A, B]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFlatMap[A, B]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFlatMap[A, B]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[B]())
	execute := true

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[B]](r.Failure())
			execute = false
		} else if r.Get().Empty() {
			execute = false
		}
	}

	if execute {

		if this.ioA != nil {

			valueA := runtime.NewWithState[A](this.state, this.ioA).UnsafeRun()
			if valueA.IsError() {
				this.value = result.OfError[*option.Option[B]](valueA.Failure())
			} else {
				optA := valueA.Get()

				if optA.NonEmpty() {
					a := optA.Get()
					this.value = runtime.NewWithState[B](this.state, this.f(a)).UnsafeRun()
				}
			}

		} else if prevEff.NonEmpty() {

			r := prevEff.Get().GetResult()

			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				runnableIO := this.f(effValue)
				runnableIO.SetState(this.state.Copy())
				runnableIO.SetDebug(this.debug)
				this.value = runtime.NewWithState[B](this.state, runnableIO).UnsafeRun()
			} else {
				util.PanicCastType("IOFlatMap",
					reflect.TypeOf(val), reflect.TypeFor[B]())
			}
		}

	}

	if this.debug {
		fmt.Println(this)
	}

	return currEff.(types.IOEffect)
}
