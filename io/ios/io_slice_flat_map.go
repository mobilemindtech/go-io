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

type IOSliceFlatMap[A any, B any] struct {
	value      *result.Result[*option.Option[[]B]]
	prevEffect types.IOEffect
	f          func(A) *types.IO[B]
	debug      bool
	state      *state.State
	debugInfo  *types.IODebugInfo
}

func NewSliceFlatMap[A any, B any](f func(A) *types.IO[B]) *IOSliceFlatMap[A, B] {
	return &IOSliceFlatMap[A, B]{f: f}
}

func (this *IOSliceFlatMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOSliceFlatMap[A, B]) SetState(st *state.State) {
	this.state = st
}

func (this *IOSliceFlatMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceFlatMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]B]()
}

func (this *IOSliceFlatMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceFlatMap[A, B]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceFlatMap[A, B]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceFlatMap[A, B]) String() string {
	return fmt.Sprintf("SliceFlatMap(%v)", this.value.String())
}

func (this *IOSliceFlatMap[A, B]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFlatMap[A, B]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFlatMap[A, B]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceFlatMap[A, B]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]B]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]B]](r.Failure())
		} else if r.Get().NonEmpty() {

			val := r.Get().GetValue()

			if effValue, ok := val.([]A); ok {

				var list []B
				for _, it := range effValue {

					runnableIO := this.f(it)
					//log.Printf("runnableIO = %v, state =%v", runnableIO, this.state)
					runnableIO.SetState(this.state.Copy())
					runnableIO.SetDebug(this.debug)
					resultIO := runtime.NewWithState[B](this.state, runnableIO).UnsafeRun()

					if resultIO.IsError() {
						this.value = result.OfError[*option.Option[[]B]](resultIO.Failure())
					} else if resultIO.Get().NonEmpty() {
						list = append(list, resultIO.Get().Get())
					}

				}
				if len(list) > 0 {
					this.value = result.OfValue(option.Some(list))
				}

			} else {
				util.PanicCastType("IOSliceFlatMap",
					reflect.TypeOf(val), reflect.TypeFor[[]B]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
