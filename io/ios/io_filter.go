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

type IOFilter[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func(A) bool
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewFilter[A any](f func(A) bool) *IOFilter[A] {
	return &IOFilter[A]{f: f}
}

func (this *IOFilter[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOFilter[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFilter[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFilter[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFilter[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOFilter[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOFilter[A]) String() string {
	return fmt.Sprintf("Filter(%v)", this.value.String())
}

func (this *IOFilter[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOFilter[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFilter[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFilter[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				if this.f(effValue) {
					this.value = result.OfValue(option.Some(effValue))
				}
			} else {
				util.PanicCastType("IOFilter",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
