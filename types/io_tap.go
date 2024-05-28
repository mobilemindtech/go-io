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

type IOTap[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func(A) bool
	debug      bool
	state      *state.State
}

func NewTap[A any](f func(A) bool) *IOTap[A] {
	return &IOTap[A]{f: f}
}

func (this *IOTap[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOTap[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOTap[A]) String() string {
	return fmt.Sprintf("Tap(%v)", this.value.String())
}

func (this *IOTap[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOTap[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOTap[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOTap[A]) UnsafeRun() IOEffect {
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
				util.PanicCastType("IOTap",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(IOEffect)
}
