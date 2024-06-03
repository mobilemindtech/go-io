package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IODebug[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	label      string
	debug      bool
}

func NewDebug[A any](label string) *IODebug[A] {
	return &IODebug[A]{label: label}
}

func (this *IODebug[A]) String() string {
	return fmt.Sprintf("Debug(%v)", this.value.String())
}

func (this *IODebug[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IODebug[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IODebug[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IODebug[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IODebug[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IODebug[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IODebug[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		log.Printf("<DEBUG>: %v - %v\n", this.label, prevEff.Get())
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().NonEmpty() {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IODebug",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	} else {
		log.Printf("<DEBUG>: %v - IO(empty)\n", this.label)
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
