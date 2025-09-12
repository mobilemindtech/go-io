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

type IOCatchAll[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	fn         func(error)
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewCatchAll[A any](f func(error)) *IOCatchAll[A] {
	return &IOCatchAll[A]{fn: f}
}

func (this *IOCatchAll[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOCatchAll[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOCatchAll[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOCatchAll[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOCatchAll[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOCatchAll[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOCatchAll[A]) String() string {
	return fmt.Sprintf("CatchAll(%v)", this.value.String())
}

func (this *IOCatchAll[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOCatchAll[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOCatchAll[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOCatchAll[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {

			this.fn(r.Failure())
			this.value = result.OfError[*option.Option[A]](r.Failure())

		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOCatchAll",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())

	}

	return currEff.(types.IOEffect)
}
