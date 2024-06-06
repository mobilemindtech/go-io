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

type IORecover[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func(error) A
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewRecover[A any](f func(error) A) *IORecover[A] {
	return &IORecover[A]{f: f}
}

func (this *IORecover[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IORecover[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IORecover[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IORecover[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IORecover[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IORecover[A]) String() string {
	return fmt.Sprintf("Recover(%v)", this.value.String())
}

func (this *IORecover[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IORecover[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IORecover[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IORecover[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfValue(option.Of(this.f(r.Failure())))
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IORecover",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())

	}

	return currEff.(types.IOEffect)
}
