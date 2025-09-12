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

type IOSliceMap[A any, B any] struct {
	value      *result.Result[*option.Option[[]B]]
	prevEffect types.IOEffect
	f          func(A) B
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewSliceMap[A any, B any](f func(A) B) *IOSliceMap[A, B] {
	return &IOSliceMap[A, B]{f: f}
}

func (this *IOSliceMap[A, B]) Lift() *types.IO[B] {
	return types.NewIO[B]().Effects(this)
}

func (this *IOSliceMap[A, B]) TypeIn() reflect.Type {
	return reflect.TypeFor[[]A]()
}

func (this *IOSliceMap[A, B]) TypeOut() reflect.Type {
	return reflect.TypeFor[[]B]()
}

func (this *IOSliceMap[A, B]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOSliceMap[A, B]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOSliceMap[A, B]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOSliceMap[A, B]) String() string {
	return fmt.Sprintf("SliceMap(%v)", this.value.String())
}

func (this *IOSliceMap[A, B]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOSliceMap[A, B]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOSliceMap[A, B]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOSliceMap[A, B]) UnsafeRun() types.IOEffect {
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
					list = append(list, this.f(it))
				}
				this.value = result.OfValue(option.Some(list))
			} else {
				util.PanicCastType("IOSliceMap",
					reflect.TypeOf(val), reflect.TypeFor[[]B]())

			}
		}
	}

	//if this.debug {
	log.Printf("%v\n", this.String())
	//}

	return currEff.(types.IOEffect)
}
