package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOOrElse[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func() *IO[A]
	debug      bool
}

func NewOrElse[A any](f func() *IO[A]) *IOOrElse[A] {
	return &IOOrElse[A]{f: f}
}

func (this *IOOrElse[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOrElse[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOOrElse[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOOrElse[A]) String() string {
	return fmt.Sprintf("OrElse(%v)", this.value.String())
}

func (this *IOOrElse[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOOrElse[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOOrElse[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOOrElse[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().Empty() {
			this.value = this.f().UnsafeRun()
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOOrElse",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
