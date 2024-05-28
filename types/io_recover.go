package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IORecover[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func(error) A
	debug      bool
	state      *state.State
}

func NewRecover[A any](f func(error) A) *IORecover[A] {
	return &IORecover[A]{f: f}
}

func (this *IORecover[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IORecover[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IORecover[A]) String() string {
	return fmt.Sprintf("Recover(%v)", this.value.String())
}

func (this *IORecover[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IORecover[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IORecover[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IORecover[A]) UnsafeRun() IOEffect {
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

	return currEff.(IOEffect)
}
