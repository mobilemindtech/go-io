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

type IOThen[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func(A) A
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewThen[A any](f func(A) A) *IOThen[A] {
	return &IOThen[A]{f: f}
}

func (this *IOThen[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOThen[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOThen[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOThen[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOThen[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOThen[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOThen[A]) String() string {
	return fmt.Sprintf("Then(%v)", this.value.String())
}

func (this *IOThen[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOThen[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOThen[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOThen[A]) UnsafeRun() types.IOEffect {
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
				this.value = result.OfValue(option.Some(this.f(effValue)))
			} else {
				util.PanicCastType("IOThen",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}
