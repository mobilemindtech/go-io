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

type IOOr[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() A
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewOr[A any](f func() A) *IOOr[A] {
	return &IOOr[A]{f: f}
}

func (this *IOOr[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOOr[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOr[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOr[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOOr[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOOr[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOOr[A]) String() string {
	return fmt.Sprintf("Or(%v)", this.value.String())
}

func (this *IOOr[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOOr[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOOr[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOOr[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().Empty() {
			this.value = result.OfValue(option.Some(this.f()))
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOOr",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}
	return currEff.(types.IOEffect)
}
