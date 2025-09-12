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

type IOAndThan[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() *types.IO[A]
	debug      bool
	debugInfo  *types.IODebugInfo
	state      *state.State
	otherIO    *types.IO[A]
}

func NewAndThan[A any](f func() *types.IO[A]) *IOAndThan[A] {
	return &IOAndThan[A]{f: f}
}

func NewAndThanIO[A any](other *types.IO[A]) *IOAndThan[A] {
	return &IOAndThan[A]{otherIO: other}
}

func (this *IOAndThan[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAndThan[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[*unit.Unit]()
}

func (this *IOAndThan[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAndThan[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAndThan[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAndThan[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAndThan[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAndThan[A]) String() string {
	return fmt.Sprintf("AndThan(%v)", this.value.String())
}

func (this *IOAndThan[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAndThan[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAndThan[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAndThan[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {

			var runnableIO types.IORunnable

			if this.otherIO != nil {
				runnableIO = this.otherIO
			} else {
				runnableIO = this.f()
			}
			runnableIO.SetPrevEffect(prevEff.Get())
			this.value = runtime.
				NewWithState[A](this.state, runnableIO).
				WithDebug(this.debug).
				UnsafeRun()
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
