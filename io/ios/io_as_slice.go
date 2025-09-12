package ios

import (
	"fmt"
	"github.com/mobilemindtech/go-io/option"
	"github.com/mobilemindtech/go-io/result"
	"github.com/mobilemindtech/go-io/types"
	"github.com/mobilemindtech/go-io/util"
	"log"
	"reflect"
)

type IOAsSlice[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewAsSliceOf[A any]() *IOAsSlice[A] {
	return &IOAsSlice[A]{}
}

func (this *IOAsSlice[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAsSlice[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOAsSlice[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOAsSlice[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAsSlice[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAsSlice[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAsSlice[A]) String() string {
	return fmt.Sprintf("AsSlice(%v)", this.value.String())
}

func (this *IOAsSlice[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAsSlice[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAsSlice[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAsSlice[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if r.Get().NonEmpty() {

			val := r.Get().GetValue()

			if effValue, ok := val.([]A); ok {
				if len(effValue) > 0 {
					this.value = result.OfValue(option.Some(effValue))
				}
			} else {
				util.PanicCastType("IOAsSlice",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
