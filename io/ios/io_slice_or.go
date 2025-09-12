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

type IOSliceOr[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func() []A
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewSliceOr[A any](f func() []A) *IOSliceOr[A] {
	return &IOSliceOr[A]{f: f}
}

func (this *IOSliceOr[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceOr[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceOr[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceOr[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceOr[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceOr[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceOr[A]) String() string {
	return fmt.Sprintf("SliceOr(%v)", this.value.String())
}

func (this *IOSliceOr[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceOr[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceOr[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceOr[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[[]A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[[]A]](r.Failure())
		} else if r.Get().IsEmpty() {
			this.value = result.OfValue(option.Some(this.f()))
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.([]A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOSliceOr",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
