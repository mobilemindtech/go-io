package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOSliceFilter[A any] struct {
	value      *result.Result[*option.Option[[]A]]
	prevEffect types.IOEffect
	f          func(A) bool
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewSliceFilter[A any](f func(A) bool) *IOSliceFilter[A] {
	return &IOSliceFilter[A]{f: f}
}

func (this *IOSliceFilter[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOSliceFilter[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceFilter[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceFilter[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceFilter[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceFilter[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceFilter[A]) String() string {
	return fmt.Sprintf("SliceFilter(%v)", this.value.String())
}

func (this *IOSliceFilter[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceFilter[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceFilter[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceFilter[A]) UnsafeRun() types.IOEffect {
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

				var list []A
				for _, it := range effValue {
					if this.f(it) {
						list = append(list, it)
					}
				}
				if len(list) > 0 {
					this.value = result.OfValue(option.Some(list))
				}

			} else {
				util.PanicCastType("IOSliceFilter",
					reflect.TypeOf(val), reflect.TypeFor[[]A]())

			}

		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
