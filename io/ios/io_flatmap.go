package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"reflect"
)

type IOFlatMap[A any, B any] struct {
	value      *result.Result[*option.Option[B]]
	prevEffect types.IOEffect
	f          func(A) types.IORunnable
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewFlatMap[A any, B any](f func(A) types.IORunnable) *IOFlatMap[A, B] {
	return &IOFlatMap[A, B]{f: f}
}

func (this *IOFlatMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOFlatMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
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

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()

		if r.IsError() {
			this.value = result.OfError[*option.Option[B]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				runnableIO := this.f(effValue)
				runnableIO.SetState(this.state.Copy())
				runnableIO.SetDebug(this.debug)
				this.value = runtime.New[B](runnableIO).UnsafeRun()
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
