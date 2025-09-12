package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/runtime"
	"github.com/mobilemindtech/go-io/state"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type IOSliceOrElse[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func() *types.IO[A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewSliceOrElse[A any](f func() *types.IO[A]) *IOSliceOrElse[A] {
	return &IOSliceOrElse[A]{f: f}
}

func (this *IOSliceOrElse[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceOrElse[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceOrElse[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceOrElse[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceOrElse[A]) String() string {
	return fmt.Sprintf("SliceOrElse(%v)", this.value.String())
}

func (this *IOSliceOrElse[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceOrElse[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceOrElse[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceOrElse[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if r.Get().IsEmpty() {
			runnableIO := this.f()
			runnableIO.SetDebug(this.debug)
			runnableIO.SetState(this.state)
			this.value = runtime.NewWithState[[]A](this.state, runnableIO).UnsafeRun()
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOSliceOrElse",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
