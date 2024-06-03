package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOForeach[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func(A)
	debug      bool
}

func NewForeach[A any](f func(A)) *IOForeach[A] {
	return &IOForeach[A]{f: f}
}

func (this *IOForeach[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOForeach[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOForeach[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOForeach[A]) String() string {
	return fmt.Sprintf("Foreach(%v)", this.value.String())
}

func (this *IOForeach[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOForeach[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOForeach[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOForeach[A]) UnsafeRun() IOEffect {
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
				this.f(effValue)
			} else {
				util.PanicCastType("IOForeach",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
