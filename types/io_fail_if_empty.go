package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOFailIfEmpty[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func() error
	debug      bool
}

func NewFailIfEmpty[A any](f func() error) *IOFailIfEmpty[A] {
	return &IOFailIfEmpty[A]{f: f}
}

func (this *IOFailIfEmpty[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIfEmpty[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOFailIfEmpty[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOFailIfEmpty[A]) String() string {
	return fmt.Sprintf("FailIfEmpty(%v)", this.value.String())
}

func (this *IOFailIfEmpty[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOFailIfEmpty[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOFailIfEmpty[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOFailIfEmpty[A]) UnsafeRun() IOEffect {
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
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOFailIfEmpty",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		} else {
			this.value = result.OfError[*option.Option[A]](this.f())
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
