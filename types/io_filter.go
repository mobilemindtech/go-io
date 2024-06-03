package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOFilter[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func(A) bool
	debug      bool
}

func NewFilter[A any](f func(A) bool) *IOFilter[A] {
	return &IOFilter[A]{f: f}
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

func (this *IOFilter[A]) String() string {
	return fmt.Sprintf("Filter(%v)", this.value.String())
}

func (this *IOFilter[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFilter[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFilter[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFilter[A]) UnsafeRun() IOEffect {
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

	return currEff.(IOEffect)
}
