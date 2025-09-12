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

type IOOrElse[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() *types.IO[A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
}

func NewOrElse[A any](f func() *types.IO[A]) *IOOrElse[A] {
	return &IOOrElse[A]{f: f}
}

func (this *IOOrElse[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOrElse[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOOrElse[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOOrElse[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOOrElse[A]) String() string {
	return fmt.Sprintf("OrElse(%v)", this.value.String())
}

func (this *IOOrElse[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOOrElse[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOOrElse[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOOrElse[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().Empty() {
			runnableIO := this.f()
			this.value = runtime.NewWithState[A](this.state, runnableIO).UnsafeRun()
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOOrElse",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
